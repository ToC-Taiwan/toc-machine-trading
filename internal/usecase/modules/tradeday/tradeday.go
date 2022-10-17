// Package tradeday package tradeday
package tradeday

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"tmt/global"
	"tmt/internal/entity"
)

// TradeDay -.
type TradeDay struct {
	holidayTimeMap map[string]bool
	tradeDayMap    map[time.Time]bool
}

// NewTradeDay -.
func NewTradeDay() *TradeDay {
	t := &TradeDay{
		holidayTimeMap: make(map[string]bool),
		tradeDayMap:    make(map[time.Time]bool),
	}

	err := t.parseHolidayFile()
	if err != nil {
		log.Panic(err)
	}

	firstDay := time.Date(global.StartTradeYear, 1, 1, 0, 0, 0, 0, time.Local)
	for {
		if firstDay.Year() > global.EndTradeYear {
			break
		}

		if firstDay.Weekday() != time.Saturday && firstDay.Weekday() != time.Sunday && !t.holidayTimeMap[firstDay.Format(global.ShortTimeLayout)] {
			t.tradeDayMap[firstDay] = true
		} else {
			t.tradeDayMap[firstDay] = false
		}

		firstDay = firstDay.AddDate(0, 0, 1)
	}
	return t
}

func (t *TradeDay) parseHolidayFile() error {
	type holidayArr struct {
		DateArr []string `json:"date_arr"`
	}

	var tmp holidayArr
	holidayFile, err := os.ReadFile("./data/holidays.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(holidayFile, &tmp); err != nil {
		return err
	}

	for _, v := range tmp.DateArr {
		t.holidayTimeMap[v] = true
	}
	return nil
}

// GetAllCalendar -.
func (t *TradeDay) GetAllCalendar() []*entity.CalendarDate {
	var calendarArr []*entity.CalendarDate
	for k, v := range t.tradeDayMap {
		calendarArr = append(calendarArr, &entity.CalendarDate{
			Date:       k,
			IsTradeDay: v,
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
	if !t.tradeDayMap[d] {
		nowTime = nowTime.AddDate(0, 0, 1)
		return t.getNextTradeDay(nowTime)
	}
	return d
}

// GetAbsNextTradeDayTime -.
func (t *TradeDay) GetAbsNextTradeDayTime(dt time.Time) time.Time {
	d := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	if !t.tradeDayMap[d] {
		dt = dt.AddDate(0, 0, 1)
		return t.GetAbsNextTradeDayTime(dt)
	}
	return d
}

// GetLastNTradeDayByDate -.
func (t *TradeDay) GetLastNTradeDayByDate(n int64, firstDay time.Time) []time.Time {
	var arr []time.Time
	for {
		if t.tradeDayMap[firstDay.AddDate(0, 0, -1)] {
			arr = append(arr, firstDay.AddDate(0, 0, -1))
		}
		if len(arr) == int(n) {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return arr
}
