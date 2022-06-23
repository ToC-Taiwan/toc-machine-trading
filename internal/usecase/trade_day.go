package usecase

import "time"

func decideTradeDay() time.Time {
	var today time.Time
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	return getNextTradeDay(today)
}

func getNextTradeDay(nowTime time.Time) time.Time {
	d := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	calendar := cc.GetCalendar()
	if !calendar[d] {
		nowTime = nowTime.AddDate(0, 0, 1)
		return getNextTradeDay(nowTime)
	}
	return d
}

func getAbsNextTradeDayTime(t time.Time) time.Time {
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	calendar := cc.GetCalendar()
	if !calendar[d] {
		t = t.AddDate(0, 0, 1)
		return getAbsNextTradeDayTime(t)
	}
	return d
}

func getLastNTradeDayByDate(n int64, firstDay time.Time) []time.Time {
	calendar := cc.GetCalendar()
	var arr []time.Time
	for {
		if calendar[firstDay.AddDate(0, 0, -1)] {
			arr = append(arr, firstDay.AddDate(0, 0, -1))
		}
		if len(arr) == int(n) {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return arr
}
