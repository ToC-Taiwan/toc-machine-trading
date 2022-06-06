package usecase

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/global"
	"toc-machine-trading/pkg/logger"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo    BasicRepo
	gRPCAPI BasicgRPCAPI
}

// NewBasic -.
func NewBasic(r *repo.BasicRepo, t *grpcapi.BasicgRPCAPI) *BasicUseCase {
	uc := &BasicUseCase{
		repo:    r,
		gRPCAPI: t,
	}

	if err := uc.importCalendarDate(context.Background()); err != nil {
		logger.Get().Panic(err)
	}

	if _, err := uc.GetAllSinopacStockAndUpdateRepo(context.Background()); err != nil {
		logger.Get().Panic(err)
	}

	return uc
}

// GetAllSinopacStockAndUpdateRepo -.
func (uc *BasicUseCase) GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error) {
	stockArr, err := uc.gRPCAPI.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	var stockDetail []*entity.Stock
	for _, v := range stockArr {
		if v.GetReference() == 0 {
			continue
		}
		stock := &entity.Stock{
			Number:    v.GetCode(),
			Name:      v.GetName(),
			Exchange:  v.GetExchange(),
			Category:  v.GetCategory(),
			DayTrade:  v.GetDayTrade() == "Yes",
			LastClose: v.GetReference(),
		}
		stockDetail = append(stockDetail, stock)

		// save to cache
		SetStockDetail(stock)
	}

	err = uc.repo.InserOrUpdatetStockArr(context.Background(), stockDetail)
	if err != nil {
		return []*entity.Stock{}, err
	}

	return stockDetail, nil
}

// GetAllRepoStock -.
func (uc *BasicUseCase) GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error) {
	data, err := uc.repo.QueryAllStock(context.Background())
	if err != nil {
		return []*entity.Stock{}, err
	}

	for _, s := range data {
		// save to cache
		SetStockDetail(s)
	}

	return data, nil
}

func parseHolidayFile() ([]string, error) {
	var holidayArr struct {
		DateArr []string `json:"date_arr"`
	}

	holidayFile, err := ioutil.ReadFile("./data/holidays.json")
	if err != nil {
		return []string{}, err
	}

	if err := json.Unmarshal(holidayFile, &holidayArr); err != nil {
		return []string{}, err
	}

	return holidayArr.DateArr, nil
}

func (uc *BasicUseCase) importCalendarDate(ctx context.Context) (err error) {
	holidayArr, err := parseHolidayFile()
	if err != nil {
		return err
	}

	holidayTimeMap := make(map[string]bool)
	for _, v := range holidayArr {
		holidayTimeMap[v] = true
	}

	firstDay := time.Date(global.StartTradeYear, 1, 1, 0, 0, 0, 0, time.Local)
	var tmp []*entity.CalendarDate
	for {
		var isTradeDay bool
		if firstDay.Year() > global.EndTradeYear {
			break
		}
		if firstDay.Weekday() != time.Saturday && firstDay.Weekday() != time.Sunday && !holidayTimeMap[firstDay.Format(global.ShortTimeLayout)] {
			isTradeDay = true
		}
		tmp = append(tmp, &entity.CalendarDate{
			Date:       firstDay,
			IsTradeDay: isTradeDay,
		})
		firstDay = firstDay.AddDate(0, 0, 1)
	}

	if err := uc.repo.InserOrUpdatetCalendarDateArr(ctx, tmp); err != nil {
		return err
	}

	return nil
}
