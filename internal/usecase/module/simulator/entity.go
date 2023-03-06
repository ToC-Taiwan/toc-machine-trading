package simulator

import (
	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
)

type SimulatorFutureTarget struct {
	Code        string
	Ticks       entity.RealTimeFutureTickArr
	TradePeriod tradeday.TradePeriod
	TradeConfig *config.TradeFuture
	Quota       *quota.Quota
}

type SimulateBalance struct {
	TotalBalance int64                 `json:"total_balance,omitempty"`
	Forward      int64                 `json:"forward,omitempty"`
	ForwardCount int64                 `json:"forward_count,omitempty"`
	ForwardOrder []*entity.FutureOrder `json:"forward_order,omitempty"`
	Reverse      int64                 `json:"reverse,omitempty"`
	ReverseCount int64                 `json:"reverse_count,omitempty"`
	ReverseOrder []*entity.FutureOrder `json:"reverse_order,omitempty"`
	Cond         *config.TradeFuture   `json:"cond,omitempty"`
}
