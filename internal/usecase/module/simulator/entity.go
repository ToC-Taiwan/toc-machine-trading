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
	TotalBalance int64
	Forward      int64
	Reverse      int64
	Cond         *config.TradeFuture
}
