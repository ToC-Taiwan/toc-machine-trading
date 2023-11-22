package dt

import (
	"time"

	"tmt/internal/config"
	"tmt/internal/entity"
)

type orderWithCfg struct {
	order           *entity.FutureOrder
	cfg             *config.TradeFuture
	maxTradeOutTime time.Time
}
