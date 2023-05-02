package usecase

import (
	"context"
	"fmt"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"

	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo     BasicRepo
	sc       BasicgRPCAPI
	cfg      *config.Config
	tradeDay *tradeday.TradeDay

	allStockDetail  []*entity.Stock
	allFutureDetail []*entity.Future
	allOptionDetail []*entity.Option
}

func (u *UseCaseBase) NewBasic() Basic {
	uc := &BasicUseCase{
		repo:     repo.NewBasic(u.pg),
		sc:       grpcapi.NewBasic(u.sc, u.cfg.Development),
		cfg:      u.cfg,
		tradeDay: tradeday.Get(),
	}

	if err := uc.importCalendarDate(context.Background()); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoStock(); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoFuture(); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoOption(); err != nil {
		logger.Fatal(err)
	}

	uc.fillBasicInfo()
	uc.saveStockFutureCache()

	return uc
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

func (uc *BasicUseCase) GetConfig() *config.Config {
	return uc.cfg
}

func (uc *BasicUseCase) fillBasicInfo() {
	tradeDay := uc.tradeDay.GetStockTradeDay().TradeDay
	openTime := 9 * time.Hour
	lastTradeDayArr := uc.tradeDay.GetLastNTradeDayByDate(2, tradeDay)

	basic := &entity.BasicInfo{
		TradeDay:           tradeDay,
		LastTradeDay:       lastTradeDayArr[0],
		BefroeLastTradeDay: lastTradeDayArr[1],

		OpenTime:       tradeDay.Add(openTime).Add(time.Duration(uc.cfg.TradeStock.HoldTimeFromOpen) * time.Second),
		TradeInEndTime: tradeDay.Add(openTime).Add(time.Duration(uc.cfg.TradeStock.TradeInEndTime) * time.Minute),
		EndTime:        tradeDay.Add(openTime).Add(time.Duration(uc.cfg.TradeStock.TotalOpenTime) * time.Minute),

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

func (uc *BasicUseCase) updateRepoStock() ([]*entity.Stock, error) {
	stockArr, err := uc.sc.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	for _, v := range stockArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
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

func (uc *BasicUseCase) updateRepoFuture() ([]*entity.Future, error) {
	futureArr, err := uc.sc.GetAllFutureDetail()
	if err != nil {
		return []*entity.Future{}, err
	}

	duplCodeMap := make(map[string]struct{})
	for _, v := range futureArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			return []*entity.Future{}, pErr
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			return []*entity.Future{}, e
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
		} else {
			logger.Warnf("Dupl future code: %s %s", v.Code, v.Name)
		}
	}

	err = uc.repo.InsertOrUpdatetFutureArr(context.Background(), uc.allFutureDetail)
	if err != nil {
		return []*entity.Future{}, err
	}
	return uc.allFutureDetail, nil
}

func (uc *BasicUseCase) updateRepoOption() ([]*entity.Option, error) {
	optionArr, err := uc.sc.GetAllOptionDetail()
	if err != nil {
		return []*entity.Option{}, err
	}

	duplCodeMap := make(map[string]struct{})
	for _, v := range optionArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			return []*entity.Option{}, pErr
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			return []*entity.Option{}, e
		}

		option := &entity.Option{
			Code:           v.GetCode(),
			Symbol:         v.GetSymbol(),
			Name:           v.GetName(),
			Category:       v.GetCategory(),
			DeliveryMonth:  v.GetDeliveryMonth(),
			DeliveryDate:   dDate.Add(810 * time.Minute),
			UnderlyingKind: v.GetUnderlyingKind(),
			StrikePrice:    v.GetStrikePrice(),
			OptionRight:    v.GetOptionRight(),
			Unit:           v.GetUnit(),
			LimitUp:        v.GetLimitUp(),
			LimitDown:      v.GetLimitDown(),
			Reference:      v.GetReference(),
			UpdateDate:     updateTime,
		}

		if _, ok := duplCodeMap[option.Code]; !ok {
			duplCodeMap[option.Code] = struct{}{}
			uc.allOptionDetail = append(uc.allOptionDetail, option)
		} else {
			logger.Warnf("Dupl option code: %s %s", v.Code, v.Name)
		}
	}

	err = uc.repo.InsertOrUpdatetOptionArr(context.Background(), uc.allOptionDetail)
	if err != nil {
		return []*entity.Option{}, err
	}
	return uc.allOptionDetail, nil
}

func (uc *BasicUseCase) importCalendarDate(ctx context.Context) error {
	if err := uc.repo.InsertOrUpdatetCalendarDateArr(ctx, uc.tradeDay.GetAllCalendar()); err != nil {
		return err
	}
	return nil
}

func (uc *BasicUseCase) saveStockFutureCache() {
	for _, s := range uc.allStockDetail {
		f, err := uc.repo.QueryFutureByLikeName(context.Background(), s.Name)
		if err != nil {
			logger.Error(err)
		}

		for _, v := range f {
			if v.Symbol == fmt.Sprintf("%sR1", v.Category) || v.Symbol == fmt.Sprintf("%sR2", v.Category) {
				continue
			}

			if time.Now().Before(v.DeliveryDate) {
				s.Future = v
				cc.SetStockDetail(s)
				break
			}
		}
	}

	for _, f := range uc.allFutureDetail {
		cc.SetFutureDetail(f)
	}
}
