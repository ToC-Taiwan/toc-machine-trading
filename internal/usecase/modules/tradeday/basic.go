package tradeday

import (
	"encoding/json"
	"os"
	"time"

	"tmt/global"
)

// TradeDay -.
type TradeDay struct {
	holidayTimeMap map[time.Time]struct{}
	tradeDayMap    map[time.Time]struct{}
}

type holidayArr struct {
	DateArr []string `json:"date_arr"`
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
