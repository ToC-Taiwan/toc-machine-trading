package simulator

import (
	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
)

type SimulatorFutureTarget struct {
	Code        string
	TradePeriod tradeday.TradePeriod

	TradeConfig *config.TradeFuture
	Quota       *quota.Quota

	Ticks []*entity.FutureHistoryTick
}

type SimulateBalance struct {
	TotalBalance int64               `json:"total_balance" yaml:"total_balance"`
	Forward      int64               `json:"forward" yaml:"forward"`
	Reverse      int64               `json:"reverse" yaml:"reverse"`
	Cond         *config.TradeFuture `json:"cond" yaml:"cond"`
}
