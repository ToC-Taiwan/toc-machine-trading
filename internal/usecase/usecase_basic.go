package usecase

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo    BasicRepo
	gRPCAPI BasicgRPCAPI
}

// NewBasic -.
func NewBasic(r BasicRepo, t BasicgRPCAPI) *BasicUseCase {
	uc := &BasicUseCase{repo: r, gRPCAPI: t}
	go func() {
		err := uc.gRPCAPI.Heartbeat()
		if err != nil {
			log.Panic(err)
		}
	}()

	if err := uc.importCalendarDate(context.Background()); err != nil {
		log.Panic(err)
	}

	if err := uc.fillBasicInfo(); err != nil {
		log.Panic(err)
	}

	if _, err := uc.GetAllSinopacStockAndUpdateRepo(context.Background()); err != nil {
		log.Panic(err)
	}
	return uc
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

	var stockDetail []*entity.Stock
	for _, v := range stockArr {
		if v.GetReference() == 0 {
			continue
		}

		updateTime, pErr := time.ParseInLocation(global.ShortSlashTimeLayout, v.GetUpdateDate(), time.Local)
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
		stockDetail = append(stockDetail, stock)

		// save stock in cache
		cc.SetStockDetail(stock)
	}

	err = uc.repo.UpdateAllStockDayTradeToNo(context.Background())
	if err != nil {
		return []*entity.Stock{}, err
	}

	err = uc.repo.InsertOrUpdatetStockArr(context.Background(), stockDetail)
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

	var result []*entity.Stock
	for _, v := range data {
		result = append(result, v)
	}
	return result, nil
}

func (uc *BasicUseCase) importCalendarDate(ctx context.Context) error {
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

	tradeDayMap := make(map[time.Time]bool)
	for _, v := range tmp {
		if v.IsTradeDay {
			tradeDayMap[v.Date] = true
		}
	}
	cc.SetCalendar(tradeDayMap)

	if err := uc.repo.InsertOrUpdatetCalendarDateArr(ctx, tmp); err != nil {
		return err
	}
	return nil
}

func (uc *BasicUseCase) fillBasicInfo() error {
	tradeDay := decideTradeDay()
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	openTime := 9 * time.Hour
	lastTradeDayArr := getLastNTradeDayByDate(2, tradeDay)
	basic := &entity.BasicInfo{
		TradeDay:           tradeDay,
		LastTradeDay:       lastTradeDayArr[0],
		BefroeLastTradeDay: lastTradeDayArr[1],

		OpenTime:        tradeDay.Add(openTime).Add(time.Duration(cfg.TradeSwitch.HoldTimeFromOpen) * time.Second),
		EndTime:         tradeDay.Add(openTime).Add(time.Duration(cfg.TradeSwitch.TotalOpenTime) * time.Minute),
		TradeInEndTime:  tradeDay.Add(openTime).Add(time.Duration(cfg.TradeSwitch.TradeInEndTime) * time.Minute),
		TradeOutEndTime: tradeDay.Add(openTime).Add(time.Duration(cfg.TradeSwitch.TradeOutEndTime) * time.Minute),

		HistoryCloseRange: getLastNTradeDayByDate(cfg.History.HistoryClosePeriod, tradeDay),
		HistoryKbarRange:  getLastNTradeDayByDate(cfg.History.HistoryKbarPeriod, tradeDay),
		HistoryTickRange:  getLastNTradeDayByDate(cfg.History.HistoryTickPeriod, tradeDay),
	}
	cc.SetBasicInfo(basic)
	return nil
}

func parseHolidayFile() ([]string, error) {
	var holidayArr struct {
		DateArr []string `json:"date_arr"`
	}

	holidayFile, err := os.ReadFile("./data/holidays.json")
	if err != nil {
		return []string{}, err
	}

	if err := json.Unmarshal(holidayFile, &holidayArr); err != nil {
		return []string{}, err
	}
	return holidayArr.DateArr, nil
}
