package simulator

import (
	"tmt/cmd/config"
)

type ConditionArr []*config.TradeFuture

func GenerateCond() []*config.TradeFuture {
	base := ConditionArr{
		&config.TradeFuture{
			Quantity:          2,
			MaxHoldTime:       30,
			TradeOutWaitTimes: 20,
			TradeTimeRange: config.TradeTimeRange{
				FirstPartDuration:  720,
				SecondPartDuration: 285,
			},

			TargetBalanceHigh: 2,
			TargetBalanceLow:  -2,
			TickInterval:      5,
			RateLimit:         5,
			RateChangeRatio:   3,
			OutInRatio:        70,
			InOutRatio:        30,
		},
	}

	base.appendTargetBalanceHigh()
	base.appendTargetBalanceLow()
	base.appendTickInterval()
	base.appendRateLimit()
	base.appendRateChangeRatio()
	base.appendOutInRatio()
	base.appendInOutRatio()

	return base
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
			if tmp.TickInterval >= 15 {
				break
			}
			tmp.TickInterval += 5
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
			if tmp.RateLimit >= 15 {
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
			if tmp.RateChangeRatio >= 15 {
				break
			}
			tmp.RateChangeRatio += 2
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
			if tmp.OutInRatio >= 90 {
				break
			}
			tmp.OutInRatio += 10
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
			if tmp.TargetBalanceHigh <= 10 {
				break
			}
			tmp.TargetBalanceHigh -= 10
			appendCond = append(appendCond, &tmp)
		}
	}

	*c = append(*c, appendCond...)
}
