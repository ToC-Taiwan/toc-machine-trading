// Package trader package trader
package trader

import (
	"time"

	"tmt/internal/entity"
	"tmt/pkg/utils"
)

// realTimeStockTickArr -.
type realTimeStockTickArr []*entity.RealTimeTick

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
		default:
			continue
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

func (c realTimeFutureTickArr) splitBySecond(last int) []realTimeFutureTickArr {
	if len(c) < 2 {
		return nil
	}

	var result []realTimeFutureTickArr
	var tmp realTimeFutureTickArr
	for i := len(c) - 2; i >= 1; i-- {
		if len(result) == last {
			return result
		}

		if c[i].TickTime.Second() == c[i-1].TickTime.Second() {
			tmp = append(tmp, c[i])
		} else {
			result = append(result, tmp)
			tmp = realTimeFutureTickArr{c[i]}
		}
	}

	return nil
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
			continue
		}
	}
	return 100 * float64(outVolume) / float64(outVolume+inVolume)
}

type realTimeKbar struct {
	open   float64
	high   float64
	low    float64
	close  float64
	volume int64
}

func (c realTimeFutureTickArr) getKbar() realTimeKbar {
	kbar := realTimeKbar{
		open:  c[0].Close,
		high:  c[0].Close,
		low:   c[0].Close,
		close: c[len(c)-1].Close,
	}

	for _, v := range c {
		kbar.volume += v.Volume

		if v.Close > kbar.high {
			kbar.high = v.Close
		}

		if v.Close < kbar.low {
			kbar.low = v.Close
		}
	}
	return kbar
}

type realTimeKbarArr []realTimeKbar

func (k realTimeKbarArr) isStable(count int) bool {
	if len(k) < count {
		return false
	}

	tmp := k[len(k)-count:]
	for i, v := range tmp {
		if i == 0 {
			continue
		}

		if v.high > tmp[i-1].high || v.low < tmp[i-1].low {
			return false
		}
	}

	return true
}

// TradeBalance -.
type TradeBalance struct {
	Count   int64 `json:"count"   yaml:"count"`
	Balance int64 `json:"balance" yaml:"balance"`
}
