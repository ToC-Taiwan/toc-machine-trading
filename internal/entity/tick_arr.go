package entity

import (
	"time"

	"tmt/pkg/utils"
)

type RealTimeStockTickArr []*RealTimeStockTick

func (c RealTimeStockTickArr) GetTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c RealTimeStockTickArr) GetOutInRatio() float64 {
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

func (c RealTimeStockTickArr) GetRSIByTickTime(preTime time.Time, count int) float64 {
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

type RealTimeFutureTickArr []*RealTimeFutureTick

func (c RealTimeFutureTickArr) AppendKbar(originalArr *kbarArr) int {
	if len(c) < 2 {
		return 0
	}

	k := kbar{
		Open:     c[0].Close,
		High:     c[0].Close,
		Low:      c[0].Close,
		Close:    c[0].Close,
		Volume:   c[0].Volume,
		KbarTime: c[0].TickTime,
	}

	for i, v := range c[1:] {
		if v.Close > k.High {
			k.High = v.Close
		}

		if v.Close < k.Low {
			k.Low = v.Close
		}

		k.Close = v.Close
		k.Volume += v.Volume

		if v.TickTime.Minute() != k.KbarTime.Minute() {
			k.KbarTime = time.Date(
				v.TickTime.Year(),
				v.TickTime.Month(),
				v.TickTime.Day(),
				v.TickTime.Hour(),
				v.TickTime.Minute(),
				0,
				0,
				v.TickTime.Location(),
			)

			if k.Close < k.Open {
				k.KbarType = KbarTypeGreen
			} else {
				k.KbarType = KbarTypeRed
			}

			*originalArr = append(*originalArr, k)
			return i + 1
		}
	}

	return 0
}

func (c RealTimeFutureTickArr) GetTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c RealTimeFutureTickArr) GetOutInRatio() float64 {
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

// func (c RealTimeFutureTickArr) splitToMultiPeriod(baseDurtion time.Duration, count int) []RealTimeFutureTickArr {
// 	var duration []time.Duration
// 	var periodArr []RealTimeFutureTickArr
// 	for i := 1; i <= count; i++ {
// 		duration = append(duration, baseDurtion*time.Duration(i))
// 		periodArr = append(periodArr, RealTimeFutureTickArr{})
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

type KbarType int

const (
	KbarTypeRed KbarType = iota + 1
	KbarTypeGreen
)

type kbar struct {
	KbarType
	KbarTime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
}

type kbarArr []kbar

func (k kbarArr) IsStable(count, limit int) bool {
	if len(k) < count {
		return false
	}

	var diff float64
	var try int
	start := k[len(k)-1].Close
	for i := len(k) - 2; i >= 0; i-- {
		diff += k[i].Close - start
		try++
		if try >= count {
			break
		}
	}

	return diff < float64(limit) && diff > 0
}
