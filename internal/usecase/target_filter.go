package usecase

import "tmt/pkg/config"

func targetFilter(close float64, volume int64, cond config.PriceVolumeLimit, isRealTime bool) bool {
	if close < cond.LimitPriceLow || close >= cond.LimitPriceHigh {
		return false
	}

	if !isRealTime && volume < cond.LimitVolume {
		return false
	}
	return true
}

func blackStockFilter(stockNum string, cond config.TargetCond) bool {
	for _, v := range cond.BlackStock {
		if v == stockNum {
			return false
		}
	}
	return true
}

func blackCatagoryFilter(category string, cond config.TargetCond) bool {
	for _, v := range cond.BlackCategory {
		if v == category {
			return false
		}
	}
	return true
}
