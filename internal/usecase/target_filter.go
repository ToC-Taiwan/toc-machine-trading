package usecase

import "toc-machine-trading/pkg/config"

func targetFilter(close float64, volume int64, cond config.PriceVolumeLimit, isRealTime bool) bool {
	if close < cond.LimitPriceLow || close >= cond.LimitPriceHigh {
		return false
	}

	if !isRealTime && volume < cond.LimitVolume {
		return false
	}
	return true
}
