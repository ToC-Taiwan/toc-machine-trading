// Package trader package trader
package trader

import (
	"time"

	"tmt/internal/entity"
	"tmt/pkg/utils"
)

// realTimeStockTickArr -.
type realTimeStockTickArr []*entity.RealTimeStockTick

func (c realTimeStockTickArr) getTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c realTimeStockTickArr) getOutInRatio() float64 {
	if len(c) == 0 {
		return 0
	}

	var outVolume, inVolume int64
	for _, v := range c {
		switch v.TickType {
		case 1:
			outVolume += v.Volume
		case 2:
			inVolume += v.Volume
		}
	}

	return 100 * float64(outVolume) / float64(outVolume+inVolume)
}

func (c realTimeStockTickArr) getRSIByTickTime(preTime time.Time, count int) float64 {
	if len(c) == 0 || preTime.IsZero() {
		return 0
	}

	var tmp []float64
	for _, v := range c {
		if v.TickTime.Equal(preTime) || v.TickTime.After(preTime) {
			tmp = append(tmp, v.Close)
		}
	}

	return utils.GenerateRSI(tmp, count)
}

// realTimeFutureTickArr -.
type realTimeFutureTickArr []*entity.RealTimeFutureTick

func (c realTimeFutureTickArr) appendKbar(originalArr *kbarArr) int {
	if len(c) < 2 {
		return 0
	}

	k := kbar{
		open:     c[0].Close,
		high:     c[0].Close,
		low:      c[0].Close,
		close:    c[0].Close,
		volume:   c[0].Volume,
		kbarTime: c[0].TickTime,
	}

	for i, v := range c[1:] {
		if v.Close > k.high {
			k.high = v.Close
		}

		if v.Close < k.low {
			k.low = v.Close
		}

		k.close = v.Close
		k.volume += v.Volume

		if v.TickTime.Minute() != k.kbarTime.Minute() {
			k.kbarTime = time.Date(
				v.TickTime.Year(),
				v.TickTime.Month(),
				v.TickTime.Day(),
				v.TickTime.Hour(),
				v.TickTime.Minute(),
				0,
				0,
				v.TickTime.Location(),
			)

			if k.close < k.open {
				k.kbarType = kbarTypeGreen
			} else {
				k.kbarType = kbarTypeRed
			}

			*originalArr = append(*originalArr, k)
			return i + 1
		}
	}

	return 0
}

func (c realTimeFutureTickArr) getTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c realTimeFutureTickArr) getOutInRatio() float64 {
	if len(c) == 0 {
		return 0
	}

	var outVolume, inVolume int64
	for _, v := range c {
		switch v.TickType {
		case 1:
			outVolume += v.Volume
		case 2:
			inVolume += v.Volume
		default:
			outVolume += v.Volume
		}
	}
	return 100 * float64(outVolume) / float64(outVolume+inVolume)
}

// func (c realTimeFutureTickArr) splitToMultiPeriod(baseDurtion time.Duration, count int) []realTimeFutureTickArr {
// 	var duration []time.Duration
// 	var periodArr []realTimeFutureTickArr
// 	for i := 1; i <= count; i++ {
// 		duration = append(duration, baseDurtion*time.Duration(i))
// 		periodArr = append(periodArr, realTimeFutureTickArr{})
// 	}

// 	startTime := c[len(c)-1].TickTime
// 	for _, v := range c {
// 		gap := startTime.Sub(v.TickTime)
// 		for i, p := range duration {
// 			if gap <= p {
// 				periodArr[i] = append(periodArr[i], v)
// 			}
// 		}
// 	}

// 	return periodArr
// }

// SimulateBalance -.
type SimulateBalance struct {
	Count   int64 `json:"count"`
	Balance int64 `json:"balance"`
}

type kbarType int

const (
	kbarTypeRed kbarType = iota + 1
	kbarTypeGreen
)

type kbar struct {
	kbarType
	kbarTime time.Time
	open     float64
	high     float64
	low      float64
	close    float64
	volume   int64
}

type kbarArr []kbar

func (k kbarArr) isStable(count, limit int) bool {
	if len(k) < count {
		return false
	}

	var diff float64
	var try int
	start := k[len(k)-1].close
	for i := len(k) - 2; i >= 0; i-- {
		diff += k[i].close - start
		try++
		if try >= count {
			break
		}
	}

	return diff < float64(limit) && diff > 0
}
