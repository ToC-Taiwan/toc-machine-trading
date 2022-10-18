// Package tradeday package tradeday
package tradeday

import (
	"time"

	"tmt/internal/entity"
)

// NewStockTradeDay -.
func NewStockTradeDay() *TradeDay {
	t := &TradeDay{
		holidayTimeMap: make(map[time.Time]struct{}),
		tradeDayMap:    make(map[time.Time]struct{}),
	}

	t.parseHolidayFile()
	t.fillTradeDay()
	return t
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

// DecideStockTradeDay -.
func (t *TradeDay) DecideStockTradeDay() time.Time {
	var today time.Time
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	return t.getNextTradeDay(today)
}

func (t *TradeDay) getNextTradeDay(nowTime time.Time) time.Time {
	d := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	if !t.isTradeDay(d) {
		nowTime = nowTime.AddDate(0, 0, 1)
		return t.getNextTradeDay(nowTime)
	}
	return d
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
