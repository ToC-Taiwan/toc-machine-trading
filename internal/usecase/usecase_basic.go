package usecase

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/modules/config"
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/pkg/common"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo    BasicRepo
	gRPCAPI BasicgRPCAPI

	cfg      *config.Config
	tradeDay *tradeday.TradeDay

	allStockDetail  []*entity.Stock
	allFutureDetail []*entity.Future
}

// NewBasic -.
func NewBasic(r BasicRepo, t BasicgRPCAPI) *BasicUseCase {
	uc := &BasicUseCase{
		repo:     r,
		gRPCAPI:  t,
		tradeDay: tradeday.NewTradeDay(),
		cfg:      config.GetConfig(),
	}

	go uc.HealthCheck()

	if err := uc.importCalendarDate(context.Background()); err != nil {
		log.Panic(err)
	}
	if _, err := uc.GetAllSinopacStockAndUpdateRepo(context.Background()); err != nil {
		log.Panic(err)
	}
	if _, err := uc.GetAllSinopacFutureAndUpdateRepo(context.Background()); err != nil {
		log.Panic(err)
	}

	uc.fillBasicInfo()
	bus.SubscribeTopic(event.TopicQueryMonitorFutureCode, uc.pubMonitorFutureCode)
	return uc
}

func (uc *BasicUseCase) HealthCheck() {
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch x := r.(type) {
			case error:
				err = x
			case string:
				err = errors.New(x)
			default:
				err = errors.New("unknown panic")
			}

			if errors.Is(err, io.EOF) {
				log.Warn("gRPC server is not ready, terminate")
			} else {
				log.Error(err)
			}
			os.Exit(0)
		}
	}()

	err := uc.gRPCAPI.Heartbeat()
	if err != nil {
		log.Panic(err)
	}
}

// TerminateSinopac -.
func (uc *BasicUseCase) TerminateSinopac(ctx context.Context) error {
	return uc.gRPCAPI.Terminate()
}

// GetAllSinopacStockAndUpdateRepo -.
func (uc *BasicUseCase) GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error) {
	stockArr, err := uc.gRPCAPI.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	for _, v := range stockArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(common.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if err != nil {
			return []*entity.Stock{}, pErr
		}

		stock := &entity.Stock{
			Number:     v.GetCode(),
			Name:       v.GetName(),
			Exchange:   v.GetExchange(),
			Category:   v.GetCategory(),
			DayTrade:   v.GetDayTrade() == entity.DayTradeYes,
			LastClose:  v.GetReference(),
			UpdateDate: updateTime,
		}
		uc.allStockDetail = append(uc.allStockDetail, stock)

		// save stock in cache
		cc.SetStockDetail(stock)
	}

	err = uc.repo.UpdateAllStockDayTradeToNo(context.Background())
	if err != nil {
		return []*entity.Stock{}, err
	}

	err = uc.repo.InsertOrUpdatetStockArr(context.Background(), uc.allStockDetail)
	if err != nil {
		return []*entity.Stock{}, err
	}
	return uc.allStockDetail, nil
}

// GetAllSinopacFutureAndUpdateRepo -.
func (uc *BasicUseCase) GetAllSinopacFutureAndUpdateRepo(ctx context.Context) ([]*entity.Future, error) {
	futureArr, err := uc.gRPCAPI.GetAllFutureDetail()
	if err != nil {
		return []*entity.Future{}, err
	}

	duplCodeMap := make(map[string]struct{})
	for _, v := range futureArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(common.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if err != nil {
			return []*entity.Future{}, pErr
		}

		dDate, e := time.ParseInLocation(common.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			return []*entity.Future{}, err
		}

		future := &entity.Future{
			Code:           v.GetCode(),
			Symbol:         v.GetSymbol(),
			Name:           v.GetName(),
			Category:       v.GetCategory(),
			DeliveryMonth:  v.GetDeliveryMonth(),
			DeliveryDate:   dDate.Add(810 * time.Minute),
			UnderlyingKind: v.GetUnderlyingKind(),
			Unit:           v.GetUnit(),
			LimitUp:        v.GetLimitUp(),
			LimitDown:      v.GetLimitDown(),
			Reference:      v.GetReference(),
			UpdateDate:     updateTime,
		}

		if _, ok := duplCodeMap[future.Code]; !ok {
			duplCodeMap[future.Code] = struct{}{}
			uc.allFutureDetail = append(uc.allFutureDetail, future)
			cc.SetFutureDetail(future)
		} else {
			log.Warnf("Dupl future code: %s %s", v.Code, v.Name)
		}
	}

	err = uc.repo.InsertOrUpdatetFutureArr(context.Background(), uc.allFutureDetail)
	if err != nil {
		return []*entity.Future{}, err
	}
	return uc.allFutureDetail, nil
}

// GetAllRepoStock -.
func (uc *BasicUseCase) GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error) {
	data, err := uc.repo.QueryAllStock(context.Background())
	if err != nil {
		return []*entity.Stock{}, err
	}

	var result []*entity.Stock
	for _, v := range data {
		result = append(result, v)
	}
	return result, nil
}

func (uc *BasicUseCase) importCalendarDate(ctx context.Context) error {
	if err := uc.repo.InsertOrUpdatetCalendarDateArr(ctx, uc.tradeDay.GetAllCalendar()); err != nil {
		return err
	}
	return nil
}

func (uc *BasicUseCase) fillBasicInfo() {
	tradeDay := uc.tradeDay.GetStockTradeDay().TradeDay
	openTime := 9 * time.Hour
	lastTradeDayArr := uc.tradeDay.GetLastNTradeDayByDate(2, tradeDay)

	basic := &entity.BasicInfo{
		TradeDay:           tradeDay,
		LastTradeDay:       lastTradeDayArr[0],
		BefroeLastTradeDay: lastTradeDayArr[1],

		OpenTime:       tradeDay.Add(openTime).Add(time.Duration(uc.cfg.StockTradeSwitch.HoldTimeFromOpen) * time.Second),
		TradeInEndTime: tradeDay.Add(openTime).Add(time.Duration(uc.cfg.StockTradeSwitch.TradeInEndTime) * time.Minute),
		EndTime:        tradeDay.Add(openTime).Add(time.Duration(uc.cfg.StockTradeSwitch.TotalOpenTime) * time.Minute),

		HistoryCloseRange: uc.tradeDay.GetLastNTradeDayByDate(uc.cfg.History.HistoryClosePeriod, tradeDay),
		HistoryKbarRange:  uc.tradeDay.GetLastNTradeDayByDate(uc.cfg.History.HistoryKbarPeriod, tradeDay),
		HistoryTickRange:  uc.tradeDay.GetLastNTradeDayByDate(uc.cfg.History.HistoryTickPeriod, tradeDay),

		AllStocks:  make(map[string]*entity.Stock),
		AllFutures: make(map[string]*entity.Future),
	}

	for _, s := range uc.allStockDetail {
		basic.AllStocks[s.Number] = s
	}

	for _, f := range uc.allFutureDetail {
		basic.AllFutures[f.Code] = f
	}

	cc.SetBasicInfo(basic)
}

func (uc *BasicUseCase) pubMonitorFutureCode() {
	futures, err := uc.repo.QueryAllMXFFuture(context.Background())
	if err != nil {
		log.Panic(err)
	}

	for _, v := range futures {
		if v.Code == "MXFR1" || v.Code == "MXFR2" {
			continue
		}

		if time.Now().Before(v.DeliveryDate) {
			bus.PublishTopicEvent(event.TopicMonitorFutureCode, v)
			return
		}
	}
}
