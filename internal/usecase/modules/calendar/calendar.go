// Package calendar package calendar
package calendar

import (
	"embed"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
)

//go:embed holidays.json
var files embed.FS

const (
	startTradeYear int = 2021
	endTradeYear   int = 2023
)

var (
	singleton *Calendar
	once      sync.Once
)

// Calendar -.
type Calendar struct {
	holidayTimeMap map[time.Time]struct{}
	tradeDayMap    map[time.Time]struct{}
}

type holidayArr struct {
	DateArr []string `json:"date_arr"`
}

func Get() *Calendar {
	if singleton == nil {
		once.Do(func() {
			t := &Calendar{
				holidayTimeMap: make(map[time.Time]struct{}),
				tradeDayMap:    make(map[time.Time]struct{}),
			}

			t.parseHolidayFile()
			t.fillTradeDay()

			singleton = t
		})
		return Get()
	}
	return singleton
}

// GetStockTradeDay -.
func (t *Calendar) GetStockTradeDay() TradePeriod {
	var nowTime time.Time
	if time.Now().Hour() >= 14 {
		nowTime = time.Now().AddDate(0, 0, 1)
	} else {
		nowTime = time.Now()
	}

	d := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	var startTime, endTime time.Time
	for {
		if t.isTradeDay(d) {
			startTime = d.Add(9 * time.Hour)
			endTime = startTime.Add(13 * time.Hour).Add(30 * time.Minute)
			break
		}
		d = d.AddDate(0, 0, 1)
	}
	return TradePeriod{startTime, endTime, d, t}
}

// GetFutureTradeDay -.
func (t *Calendar) GetFutureTradeDay() TradePeriod {
	var nowTime time.Time
	if time.Now().Hour() >= 14 {
		nowTime = time.Now().AddDate(0, 0, 1)
	} else {
		nowTime = time.Now()
	}

	var startTime, endTime, tradeDay time.Time
	d := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	for {
		if t.isTradeDay(d) {
			tradeDay = d
			endTime = d.Add(13 * time.Hour).Add(45 * time.Minute)
			break
		}
		d = d.AddDate(0, 0, 1)
	}

	d = d.AddDate(0, 0, -1)
	for {
		if t.isTradeDay(d) {
			startTime = d.Add(15 * time.Hour)
			break
		}
		d = d.AddDate(0, 0, -1)
	}

	return TradePeriod{startTime, endTime, tradeDay, t}
}

// GetStockTradePeriodByDate -.
func (t *Calendar) GetStockTradePeriodByDate(date string) (TradePeriod, error) {
	d, err := time.ParseInLocation(entity.ShortTimeLayout, date, time.Local)
	if err != nil {
		return TradePeriod{}, err
	}

	var startTime, endTime time.Time
	if t.isTradeDay(d) {
		startTime = d.Add(9 * time.Hour)
		endTime = startTime.Add(13 * time.Hour).Add(30 * time.Minute)
	} else {
		return TradePeriod{}, errors.New("not trade day")
	}
	return TradePeriod{startTime, endTime, d, t}, nil
}

// GetFutureTradePeriodByDate -.
func (t *Calendar) GetFutureTradePeriodByDate(date string) (TradePeriod, error) {
	d, err := time.ParseInLocation(entity.ShortTimeLayout, date, time.Local)
	if err != nil {
		return TradePeriod{}, err
	}

	var startTime, endTime, tradeDay time.Time
	if t.isTradeDay(d) {
		tradeDay = d
		endTime = d.Add(13 * time.Hour).Add(45 * time.Minute)
	} else {
		return TradePeriod{}, errors.New("not trade day")
	}

	d = d.AddDate(0, 0, -1)
	for {
		if t.isTradeDay(d) {
			startTime = d.Add(15 * time.Hour)
			break
		}
		d = d.AddDate(0, 0, -1)
	}

	return TradePeriod{startTime, endTime, tradeDay, t}, nil
}

// GetLastNFutureTradeDay -.
func (t *Calendar) GetLastNFutureTradeDay(count int) []TradePeriod {
	firstDay := t.GetFutureTradeDay()
	d := firstDay.TradeDay.AddDate(0, 0, -1)

	var tradePeriodArr []TradePeriod
	for {
		if len(tradePeriodArr) == count {
			break
		}

		var startTime, endTime, tradeDay time.Time
		for {
			if t.isTradeDay(d) {
				tradeDay = d
				endTime = d.Add(13 * time.Hour).Add(45 * time.Minute)
				break
			}
			d = d.AddDate(0, 0, -1)
		}

		d = d.AddDate(0, 0, -1)
		for {
			if t.isTradeDay(d) {
				startTime = d.Add(15 * time.Hour)
				break
			}
			d = d.AddDate(0, 0, -1)
		}

		tradePeriodArr = append(tradePeriodArr, TradePeriod{startTime, endTime, tradeDay, t})
	}

	return tradePeriodArr
}

func (t *Calendar) parseHolidayFile() {
	tmp := holidayArr{}
	content, err := files.ReadFile("holidays.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(content, &tmp); err != nil {
		panic(err)
	}

	for _, v := range tmp.DateArr {
		tm, err := time.ParseInLocation(entity.ShortTimeLayout, v, time.Local)
		if err != nil {
			panic(err)
		}

		t.holidayTimeMap[tm] = struct{}{}
	}
}

func (t *Calendar) fillTradeDay() {
	tm := time.Date(startTradeYear, 1, 1, 0, 0, 0, 0, time.Local)
	for {
		if tm.Year() > endTradeYear {
			break
		}

		if tm.Weekday() != time.Saturday && tm.Weekday() != time.Sunday && !t.isHoliday(tm) {
			t.tradeDayMap[tm] = struct{}{}
		}

		tm = tm.AddDate(0, 0, 1)
	}
}

func (t *Calendar) isHoliday(date time.Time) bool {
	if _, ok := t.holidayTimeMap[date]; ok {
		return true
	}
	return false
}

func (t *Calendar) isTradeDay(date time.Time) bool {
	if _, ok := t.tradeDayMap[date]; ok {
		return true
	}
	return false
}

// GetAllCalendar -.
func (t *Calendar) GetAllCalendar() []*entity.CalendarDate {
	var calendarArr []*entity.CalendarDate
	for k := range t.tradeDayMap {
		calendarArr = append(calendarArr, &entity.CalendarDate{
			Date:       k,
			IsTradeDay: t.isTradeDay(k),
		})
	}
	return calendarArr
}

// GetAbsNextTradeDayTime -.
func (t *Calendar) GetAbsNextTradeDayTime(dt time.Time) time.Time {
	d := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	if !t.isTradeDay(d) {
		dt = dt.AddDate(0, 0, 1)
		return t.GetAbsNextTradeDayTime(dt)
	}
	return d
}

// GetLastNTradeDayByDate -.
func (t *Calendar) GetLastNTradeDayByDate(n int64, firstDay time.Time) []time.Time {
	var arr []time.Time
	for {
		if t.isTradeDay(firstDay.AddDate(0, 0, -1)) {
			arr = append(arr, firstDay.AddDate(0, 0, -1))
		}
		if len(arr) == int(n) {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return arr
}

func (t *Calendar) GetLastNStockTradeDay(n int64) []time.Time {
	firstDay := t.GetStockTradeDay().TradeDay
	var arr []time.Time
	for {
		if t.isTradeDay(firstDay.AddDate(0, 0, -1)) {
			arr = append(arr, firstDay.AddDate(0, 0, -1))
		}
		if len(arr) == int(n) {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return arr
}

// TradePeriod -.
type TradePeriod struct {
	StartTime time.Time
	EndTime   time.Time
	TradeDay  time.Time
	base      *Calendar
}

func (tp *TradePeriod) ToTimeRange(firstMinute, secondMinute int64) [][]time.Time {
	var timeRange [][]time.Time
	timeRange = append(timeRange, []time.Time{
		tp.StartTime,
		tp.StartTime.Add(time.Duration(firstMinute) * time.Minute),
	})
	timeRange = append(timeRange, []time.Time{
		tp.EndTime.Add(-300 * time.Minute),
		tp.EndTime.Add(-300 * time.Minute).Add(time.Duration(secondMinute) * time.Minute),
	})
	return timeRange
}

func (tp *TradePeriod) IsStockMarketOpenNow() bool {
	if time.Now().After(tp.StartTime) && time.Now().Before(tp.EndTime) {
		return true
	}
	return false
}

func (tp *TradePeriod) IsFutureMarketOpenNow() bool {
	firstEndTime := tp.StartTime.Add(14 * time.Hour)
	secondStartTime := tp.EndTime.Add(-5 * time.Hour)

	now := time.Now()
	if now.After(tp.StartTime) && now.Before(firstEndTime) {
		return true
	}
	if now.After(secondStartTime) && now.Before(tp.EndTime) {
		return true
	}
	return false
}

// ToStartEndArray -.
func (tp *TradePeriod) ToStartEndArray() []time.Time {
	return []time.Time{tp.StartTime, tp.EndTime}
}

// GetLastFutureTradePeriod -.
func (tp *TradePeriod) GetLastFutureTradePeriod() TradePeriod {
	firstDay := tp
	d := firstDay.TradeDay.AddDate(0, 0, -1)

	var startTime, endTime, tradeDay time.Time
	for {
		if tp.base.isTradeDay(d) {
			tradeDay = d
			endTime = d.Add(13 * time.Hour).Add(45 * time.Minute)
			break
		}
		d = d.AddDate(0, 0, -1)
	}

	d = d.AddDate(0, 0, -1)
	for {
		if tp.base.isTradeDay(d) {
			startTime = d.Add(15 * time.Hour)
			break
		}
		d = d.AddDate(0, 0, -1)
	}
	return TradePeriod{startTime, endTime, tradeDay, tp.base}
}
