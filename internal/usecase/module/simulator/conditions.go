package simulator

import (
	"tmt/cmd/config"
)

type ConditionArr []*config.TradeFuture

func GenerateCond() []*config.TradeFuture {
	base := ConditionArr{
		&config.TradeFuture{
			Quantity:    2,
			MaxHoldTime: 15,
			TradeTimeRange: config.TradeTimeRange{
				FirstPartDuration:  720,
				SecondPartDuration: 285,
			},

			TradeOutWaitTimes: 20,
			TargetBalanceHigh: 2,
			TargetBalanceLow:  -2,
			TickInterval:      8,
			RateLimit:         8,
			RateChangeRatio:   1,
			OutInRatio:        55,
			InOutRatio:        55,
		},
	}

	base.appendTargetBalanceHigh()
	base.appendTargetBalanceLow()
	base.appendTickInterval()
	base.appendRateLimit()
	base.appendRateChangeRatio()
	base.appendOutInRatio()
	base.appendInOutRatio()
	base.appendTradeOutTimes()

	return base
}

func (c *ConditionArr) appendTradeOutTimes() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.TradeOutWaitTimes >= 60 {
				break
			}
			tmp.TradeOutWaitTimes += 5
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendTargetBalanceHigh() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.TargetBalanceHigh >= 5 {
				break
			}
			tmp.TargetBalanceHigh++
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendTargetBalanceLow() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.TargetBalanceLow <= -5 {
				break
			}
			tmp.TargetBalanceLow--
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendTickInterval() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.TickInterval >= 14 {
				break
			}
			tmp.TickInterval += 2
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendRateLimit() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.RateLimit >= 14 {
				break
			}
			tmp.RateLimit += 2
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendRateChangeRatio() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.RateChangeRatio >= 10 {
				break
			}
			tmp.RateChangeRatio++
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendOutInRatio() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.OutInRatio >= 75 {
				break
			}
			tmp.OutInRatio += 5
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}

func (c *ConditionArr) appendInOutRatio() {
	var appendCond []*config.TradeFuture
	for _, cond := range *c {
		tmp := *cond
		for {
			if tmp.InOutRatio >= 75 {
				break
			}
			tmp.InOutRatio += 5
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}
