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

func (c realTimeFutureTickArr) getActionByPeriodOutInRatioTrend(count, unit int) entity.OrderAction {
	if len(c) == 0 || len(c) < count*unit {
		return entity.ActionNone
	}

	period := []realTimeFutureTickArr{}
	single := realTimeFutureTickArr{}
	for i := len(c) - 1; i >= 0; i-- {
		v := c[i]
		single = append(single, v)

		if len(single)%unit == 0 {
			period = append(period, single)
		}

		if len(period) == count {
			break
		}
	}

	arr := outInRatioArr{}
	for _, v := range period {
		arr = append(arr, v.getOutInRatio())
	}

	return arr.getAction(count)
}

type outInRatioArr []float64

func (o outInRatioArr) getAction(count int) entity.OrderAction {
	if count%2 != 1 {
		return entity.ActionNone
	}

	// TODO: out in ratio should come from config
	middle := count / 2
	switch {
	case o[0] > o[middle] && o[middle] > o[count-1] && o[0] < 40:
		return entity.ActionBuy
	case o[0] < o[middle] && o[middle] < o[count-1] && o[0] > 60:
		return entity.ActionSellFirst
	}
	return entity.ActionNone
}

// func (c realTimeFutureTickArr) splitBySecond(last int) []realTimeFutureTickArr {
// 	if len(c) < 2 {
// 		return nil
// 	}

// 	var result []realTimeFutureTickArr
// 	var tmp realTimeFutureTickArr
// 	for i := len(c) - 2; i >= 1; i-- {
// 		if len(result) == last {
// 			return result
// 		}

// 		if c[i].TickTime.Second() == c[i-1].TickTime.Second() {
// 			tmp = append(tmp, c[i])
// 		} else {
// 			result = append(result, tmp)
// 			tmp = realTimeFutureTickArr{c[i]}
// 		}
// 	}

// 	return nil
// }

// func (c realTimeFutureTickArr) getTotalVolume() int64 {
// 	var volume int64
// 	for _, v := range c {
// 		volume += v.Volume
// 	}
// 	return volume
// }

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

// TradeBalance -.
type TradeBalance struct {
	Count   int64 `json:"count"   yaml:"count"`
	Balance int64 `json:"balance" yaml:"balance"`
}
