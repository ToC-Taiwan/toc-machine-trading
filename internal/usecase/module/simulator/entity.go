package simulator

import (
	"fmt"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
)

type SimulatorFutureTarget struct {
	Code           string
	Ticks          entity.RealTimeFutureTickArr
	TradePeriod    tradeday.TradePeriod
	Quota          *quota.Quota
	TradeConfigArr []*config.TradeFuture
}

type SimulateBalance struct {
	TotalBalance int64                 `json:"total_balance"`
	Forward      int64                 `json:"forward"`
	ForwardCount int64                 `json:"forward_count"`
	ForwardOrder []*entity.FutureOrder `json:"forward_order"`
	Reverse      int64                 `json:"reverse"`
	ReverseCount int64                 `json:"reverse_count"`
	ReverseOrder []*entity.FutureOrder `json:"reverse_order"`
	Cond         *config.TradeFuture   `json:"cond"`
}

func (s *SimulateBalance) String() string {
	return fmt.Sprintf(`
:speech_balloon: Cond: %s
Balance: %d
Forward: %d
ForwardCount: %d
Reverse: %d
ReverseCount: %d
	`, s.Cond.String(), s.TotalBalance, s.Forward, s.ForwardCount, s.Reverse, s.ReverseCount)
}
