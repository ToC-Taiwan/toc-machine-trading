package usecase

import (
	"context"
	"time"

	"tmt/internal/config"

	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/grpc"
	"tmt/internal/usecase/modules/calendar"
	"tmt/internal/usecase/repo"
	"tmt/pkg/log"
	"tmt/pkg/utils"
)

type BasicUseCase struct {
	repo     BasicRepo
	sc       BasicgRPCAPI
	cfg      *config.Config
	tradeDay *calendar.Calendar

	allStockDetail  []*entity.Stock
	allFutureDetail []*entity.Future
	allOptionDetail []*entity.Option

	logger *log.Log
	cc     *cache.Cache
}

func NewBasic() Basic {
	cfg := config.Get()
	uc := &BasicUseCase{
		repo:     repo.NewBasic(cfg.GetPostgresPool()),
		sc:       grpc.NewBasic(cfg.GetSinopacPool()),
		cfg:      cfg,
		tradeDay: calendar.Get(),
		logger:   log.Get(),
		cc:       cache.Get(),
	}

	uc.loginAll()
	uc.checkgRPCHealth()

	if err := uc.importCalendarDate(); err != nil {
		uc.logger.Fatal(err)
	}

	if err := uc.updateRepoStock(); err != nil {
		uc.logger.Fatal(err)
	}

	if err := uc.updateRepoFuture(); err != nil {
		uc.logger.Fatal(err)
	}

	if err := uc.updateRepoOption(); err != nil {
		uc.logger.Fatal(err)
	}

	return uc
}

func (uc *BasicUseCase) checkgRPCHealth() {
	go func() {
		if err := uc.sc.CreateLongConnection(); err != nil {
			uc.logger.Fatalf("sinopac CreateLongConnection error: %v", err)
		}
	}()
}

func (uc *BasicUseCase) loginAll() {
	if err := uc.sc.Login(); err != nil {
		uc.logger.Fatal(err)
	}
}

func (uc *BasicUseCase) importCalendarDate() error {
	return uc.repo.InsertOrUpdatetCalendarDateArr(context.Background(), uc.tradeDay.GetAllCalendar())
}

func (uc *BasicUseCase) GetShioajiUsage() (*entity.ShioajiUsage, error) {
	usage, err := uc.sc.CheckUsage()
	if err != nil {
		return nil, err
	}
	return &entity.ShioajiUsage{
		Connections:          int(usage.GetConnections()),
		TrafficUsage:         utils.Round(float64(usage.GetBytes())/1024/1024, 2),
		TrafficUsagePercents: utils.Round(100*(1-float64(usage.GetRemainingBytes())/float64(usage.GetLimitBytes())), 2),
	}, nil
}

func (uc *BasicUseCase) updateRepoStock() error {
	stockArr, err := uc.sc.GetAllStockDetail()
	if err != nil {
		return err
	}

	for _, v := range stockArr {
		if v.GetCode() == "001" {
			continue
		}

		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(entity.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			uc.logger.Warnf("stock %s update date parse error: %s", v.GetCode(), pErr.Error())
			continue
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
		uc.cc.SetStockDetail(stock)
	}

	err = uc.repo.UpdateAllStockDayTradeToNo(context.Background())
	if err != nil {
		return err
	}

	return uc.repo.InsertOrUpdatetStockArr(context.Background(), uc.allStockDetail)
}

func (uc *BasicUseCase) updateRepoFuture() error {
	futureArr, err := uc.sc.GetAllFutureDetail()
	if err != nil {
		return err
	}

	duplCodeMap := make(map[string]struct{})
	for _, v := range futureArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(entity.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			uc.logger.Warnf("future %s update date parse error: %s", v.GetCode(), pErr.Error())
			continue
		}

		dDate, e := time.ParseInLocation(entity.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			uc.logger.Warnf("future %s delivery date parse error: %s", v.GetCode(), e.Error())
			continue
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
			uc.cc.SetFutureDetail(future)
		}
	}

	return uc.repo.InsertOrUpdatetFutureArr(context.Background(), uc.allFutureDetail)
}

func (uc *BasicUseCase) updateRepoOption() error {
	optionArr, err := uc.sc.GetAllOptionDetail()
	if err != nil {
		return err
	}

	duplCodeMap := make(map[string]struct{})
	for _, v := range optionArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(entity.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			uc.logger.Warnf("option %s update date parse error: %s", v.GetCode(), pErr.Error())
			continue
		}

		dDate, e := time.ParseInLocation(entity.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			uc.logger.Warnf("option %s delivery date parse error: %s", v.GetCode(), e.Error())
			continue
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
		}
	}

	return uc.repo.InsertOrUpdatetOptionArr(context.Background(), uc.allOptionDetail)
}

func (uc *BasicUseCase) GetStockDetail(stockNum string) *entity.Stock {
	return uc.cc.GetStockDetail(stockNum)
}

func (uc *BasicUseCase) GetFutureDetail(code string) *entity.Future {
	return uc.cc.GetFutureDetail(code)
}

func (uc *BasicUseCase) GetConfig() *config.Config {
	return uc.cfg
}
