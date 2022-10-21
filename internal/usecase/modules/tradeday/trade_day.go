// Package tradeday package tradeday
package tradeday

import (
	"encoding/json"
	"os"
	"time"

	"tmt/global"
	"tmt/internal/entity"
)

// TradeDay -.
type TradeDay struct {
	holidayTimeMap map[time.Time]struct{}
	tradeDayMap    map[time.Time]struct{}
}

type holidayArr struct {
	DateArr []string `json:"date_arr"`
}

// NewTradeDay -.
func NewTradeDay() *TradeDay {
	t := &TradeDay{
		holidayTimeMap: make(map[time.Time]struct{}),
		tradeDayMap:    make(map[time.Time]struct{}),
	}

	t.parseHolidayFile()
	t.fillTradeDay()
	return t
}

// GetStockTradeDay -.
func (t *TradeDay) GetStockTradeDay() TradePeriod {
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
	return TradePeriod{startTime, endTime, d}
}

// GetFutureTradeDay -.
func (t *TradeDay) GetFutureTradeDay() TradePeriod {
	var nowTime time.Time
	if time.Now().Hour() >= 14 {
		nowTime = time.Now()
	} else {
		nowTime = time.Now().AddDate(0, 0, -1)
	}

	d := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	var startTime, endTime, tradeDay time.Time
	for {
		if startTime.IsZero() {
			if t.isTradeDay(d) {
				tradeDay = d
				startTime = d
				endTime = d.AddDate(0, 0, 1)
			} else {
				d = d.AddDate(0, 0, -1)
			}
		}

		if t.isTradeDay(endTime) {
			startTime = startTime.Add(15 * time.Hour)
			endTime = endTime.Add(13 * time.Hour).Add(45 * time.Minute)
			break
		}

		endTime = endTime.AddDate(0, 0, 1)
	}
	return TradePeriod{startTime, endTime, tradeDay}
}

func (t *TradeDay) parseHolidayFile() {
	tmp := holidayArr{}
	holidayFile, err := os.ReadFile("./data/holidays.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(holidayFile, &tmp); err != nil {
		panic(err)
	}

	for _, v := range tmp.DateArr {
		tm, err := time.ParseInLocation(global.ShortTimeLayout, v, time.Local)
		if err != nil {
			panic(err)
		}

		t.holidayTimeMap[tm] = struct{}{}
	}
}

func (t *TradeDay) fillTradeDay() {
	tm := time.Date(global.StartTradeYear, 1, 1, 0, 0, 0, 0, time.Local)
	for {
		if tm.Year() > global.EndTradeYear {
			break
		}

		if tm.Weekday() != time.Saturday && tm.Weekday() != time.Sunday && !t.isHoliday(tm) {
			t.tradeDayMap[tm] = struct{}{}
		}

		tm = tm.AddDate(0, 0, 1)
	}
}

func (t *TradeDay) isHoliday(date time.Time) bool {
	if _, ok := t.holidayTimeMap[date]; ok {
		return true
	}
	return false
}

func (t *TradeDay) isTradeDay(date time.Time) bool {
	if _, ok := t.tradeDayMap[date]; ok {
		return true
	}
	return false
}

// GetAllCalendar -.
func (t *TradeDay) GetAllCalendar() []*entity.CalendarDate {
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
func (t *TradeDay) GetAbsNextTradeDayTime(dt time.Time) time.Time {
	d := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	if !t.isTradeDay(d) {
		dt = dt.AddDate(0, 0, 1)
		return t.GetAbsNextTradeDayTime(dt)
	}
	return d
}

// GetLastNTradeDayByDate -.
func (t *TradeDay) GetLastNTradeDayByDate(n int64, firstDay time.Time) []time.Time {
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
}

// ToStartEndArray -.
func (tp *TradePeriod) ToStartEndArray() []time.Time {
	return []time.Time{tp.StartTime, tp.EndTime}
}
