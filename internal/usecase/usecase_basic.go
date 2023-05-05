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

	if err := uc.importCalendarDate(); err != nil {
		logger.Fatal(err)
	}

	if err := uc.updateRepoStock(); err != nil {
		logger.Fatal(err)
	}

	if err := uc.updateRepoFuture(); err != nil {
		logger.Fatal(err)
	}

	if err := uc.updateRepoOption(); err != nil {
		logger.Fatal(err)
	}

	uc.saveStockFutureCache()

	return uc
}

func (uc *BasicUseCase) importCalendarDate() error {
	return uc.repo.InsertOrUpdatetCalendarDateArr(context.Background(), uc.tradeDay.GetAllCalendar())
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
			logger.Warnf("stock %s reference is 0", v.GetCode())
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			return pErr
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
			logger.Warnf("future %s reference is 0", v.GetCode())
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			return pErr
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			return e
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
			logger.Warnf("option %s reference is 0", v.GetCode())
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
		if pErr != nil {
			return pErr
		}

		dDate, e := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetDeliveryDate(), time.Local)
		if e != nil {
			return e
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

	return uc.repo.InsertOrUpdatetOptionArr(context.Background(), uc.allOptionDetail)
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
