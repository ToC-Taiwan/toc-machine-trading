package dt

import (
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
)

type orderWithCfg struct {
	order            *entity.FutureOrder
	cfg              *config.TradeFuture
	lastTradeOutTime time.Time
}
