package usecase

import (
	"context"
	"os"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"

	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
	"tmt/pkg/common"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo BasicRepo

	sc BasicgRPCAPI
	fg BasicgRPCAPI

	cfg      *config.Config
	tradeDay *tradeday.TradeDay

	allStockDetail  []*entity.Stock
	allFutureDetail []*entity.Future

	// stockTradeInSwitch  bool
	// futureTradeInSwitch bool
}

func (u *UseCaseBase) NewBasic() Basic {
	uc := &BasicUseCase{
		repo:     repo.NewBasic(u.pg),
		sc:       grpcapi.NewBasic(u.sc),
		fg:       grpcapi.NewBasic(u.fg),
		tradeDay: tradeday.Get(),
		cfg:      u.cfg,
	}

	go uc.healthCheckforSinopac()
	go uc.healthCheckforFugle()

	if err := uc.importCalendarDate(context.Background()); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoStock(); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoFuture(); err != nil {
		logger.Fatal(err)
	}

	uc.fillBasicInfo()

	// go uc.checkStockTradeSwitch()
	// go uc.checkFutureTradeSwitch()

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

// TerminateSinopac -.
func (uc *BasicUseCase) TerminateSinopac() error {
	return uc.sc.Terminate()
}

// TerminateFugle -.
func (uc *BasicUseCase) TerminateFugle() error {
	return uc.fg.Terminate()
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

func (uc *BasicUseCase) healthCheckforSinopac() {
	err := uc.sc.Heartbeat()
	if err != nil {
		logger.Warn("sinopac healthcheck fail, terminate")
		os.Exit(0)
	}
}

func (uc *BasicUseCase) healthCheckforFugle() {
	err := uc.fg.Heartbeat()
	if err != nil {
		logger.Warn("fugle healthcheck fail, terminate")
		os.Exit(0)
	}
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
			logger.Warnf("Dupl future code: %s %s", v.Code, v.Name)
		}
	}

	err = uc.repo.InsertOrUpdatetFutureArr(context.Background(), uc.allFutureDetail)
	if err != nil {
		return []*entity.Future{}, err
	}
	return uc.allFutureDetail, nil
}

func (uc *BasicUseCase) importCalendarDate(ctx context.Context) error {
	if err := uc.repo.InsertOrUpdatetCalendarDateArr(ctx, uc.tradeDay.GetAllCalendar()); err != nil {
		return err
	}
	return nil
}

// func (uc *BasicUseCase) checkStockTradeSwitch() {
// 	if !uc.cfg.TradeStock.AllowTrade {
// 		return
// 	}

// 	openTime := uc.basic.OpenTime
// 	tradeInEndTime := uc.basic.TradeInEndTime

// 	for range time.NewTicker(2500 * time.Millisecond).C {
// 		now := time.Now()
// 		var tempSwitch bool
// 		switch {
// 		case now.Before(openTime) || now.After(tradeInEndTime):
// 			tempSwitch = false
// 		case now.After(openTime) && now.Before(tradeInEndTime):
// 			tempSwitch = true
// 		}

// 		if uc.stockTradeInSwitch != tempSwitch {
// 			uc.stockTradeInSwitch = tempSwitch
// 			bus.PublishTopicEvent(topic.TopicUpdateStockTradeSwitch, uc.stockTradeInSwitch)
// 		}
// 	}
// }

// func (uc *BasicUseCase) checkFutureTradeSwitch() {
// 	if !uc.cfg.TradeFuture.AllowTrade {
// 		return
// 	}

// 	futureTradeDay := uc.tradeDay.GetFutureTradeDay()
// 	timeRange := [][]time.Time{}
// 	firstStart := futureTradeDay.StartTime
// 	secondStart := futureTradeDay.EndTime.Add(-300 * time.Minute)

// 	timeRange = append(timeRange, []time.Time{firstStart, firstStart.Add(time.Duration(uc.cfg.TradeFuture.TradeTimeRange.FirstPartDuration) * time.Minute)})
// 	timeRange = append(timeRange, []time.Time{secondStart, secondStart.Add(time.Duration(uc.cfg.TradeFuture.TradeTimeRange.SecondPartDuration) * time.Minute)})

// 	for range time.NewTicker(2500 * time.Millisecond).C {
// 		now := time.Now()
// 		var tempSwitch bool
// 		for _, rangeTime := range timeRange {
// 			if now.After(rangeTime[0]) && now.Before(rangeTime[1]) {
// 				tempSwitch = true
// 			}
// 		}

// 		if uc.futureTradeInSwitch != tempSwitch {
// 			uc.futureTradeInSwitch = tempSwitch
// 			bus.PublishTopicEvent(topic.TopicUpdateFutureTradeSwitch, uc.futureTradeInSwitch)
// 		}
// 	}
// }
