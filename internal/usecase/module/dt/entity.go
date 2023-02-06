package dt

import (
	"tmt/cmd/config"
	"tmt/internal/entity"
)

type orderWithCfg struct {
	order *entity.FutureOrder
	cfg   *config.TradeFuture
}
