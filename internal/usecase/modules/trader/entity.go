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

func (c realTimeFutureTickArr) splitBySecond() []realTimeFutureTickArr {
	if len(c) < 2 {
		return nil
	}

	var result []realTimeFutureTickArr
	var tmp realTimeFutureTickArr
	for i, tick := range c {
		if i == len(c)-1 {
			result = append(result, tmp)
			break
		}

		if tick.TickTime.Second() == c[i+1].TickTime.Second() {
			tmp = append(tmp, tick)
		} else {
			result = append(result, tmp)
			tmp = realTimeFutureTickArr{tick}
		}
	}

	return result
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