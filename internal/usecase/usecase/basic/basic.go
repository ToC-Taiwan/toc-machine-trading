// Package basic package basic
package basic

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
	"tmt/internal/utils"
	"tmt/pkg/log"
)

type BasicUseCase struct {
	repo     BasicRepo
	sc       BasicgRPCAPI
	fugle    BasicgRPCAPI
	cfg      *config.Config
	tradeDay *tradeday.TradeDay

	allStockDetail  []*entity.Stock
	allFutureDetail []*entity.Future
	allOptionDetail []*entity.Option

	logger *log.Log
	cc     *cache.Cache
}

func NewBasic() Basic {
	cfg := config.Get()
	return &BasicUseCase{
		repo:     repo.NewBasic(cfg.GetPostgresPool()),
		sc:       grpcapi.NewBasic(cfg.GetSinopacPool()),
		fugle:    grpcapi.NewBasic(cfg.GetFuglePool()),
		cfg:      cfg,
		tradeDay: tradeday.Get(),
	}
}

func (uc *BasicUseCase) Init(logger *log.Log, cc *cache.Cache) Basic {
	uc.logger = logger
	uc.cc = cc

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

	uc.saveStockFutureCache()
	return uc
}

func (uc *BasicUseCase) checkgRPCHealth() {
	go func() {
		errChan := make(chan error)
		go func(errChan chan error) {
			if err := uc.sc.CreateLongConnection(); err != nil {
				errChan <- errors.New("sinopac CreateLongConnection error")
			}
		}(errChan)
		go func(errChan chan error) {
			if err := uc.fugle.CreateLongConnection(); err != nil {
				errChan <- errors.New("fugle CreateLongConnection error")
			}
		}(errChan)
		err := <-errChan
		uc.logger.Fatal(err)
	}()
}

func (uc *BasicUseCase) loginAll() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := uc.sc.Login(); err != nil {
			uc.logger.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := uc.fugle.Login(); err != nil {
			uc.logger.Fatal(err)
		}
	}()
	wg.Wait()
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
		Connections:  int(usage.GetConnections()),
		TrafficUsage: utils.Round(float64(usage.GetBytes())/1024/1024, 2),
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

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
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

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			uc.logger.Warnf("future %s update date parse error: %s", v.GetCode(), pErr.Error())
			continue
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
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

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			uc.logger.Warnf("option %s update date parse error: %s", v.GetCode(), pErr.Error())
			continue
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
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

func (uc *BasicUseCase) saveStockFutureCache() {
	for _, s := range uc.allStockDetail {
		f, err := uc.repo.QueryFutureByLikeName(context.Background(), s.Name)
		if err != nil {
			uc.logger.Error(err)
		}

		for _, v := range f {
			if v.Symbol == fmt.Sprintf("%sR1", v.Category) || v.Symbol == fmt.Sprintf("%sR2", v.Category) {
				continue
			}

			if time.Now().Before(v.DeliveryDate) {
				s.Future = v
				uc.cc.SetStockDetail(s)
				break
			}
		}
	}

	for _, f := range uc.allFutureDetail {
		uc.cc.SetFutureDetail(f)
	}
}

func (uc *BasicUseCase) GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error) {
	data, err := uc.repo.QueryAllStock(ctx)
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