// Package trader package trader
package trader

import (
	"tmt/internal/usecase/modules/cache"
	"tmt/internal/usecase/modules/event"
	"tmt/pkg/log"
)

var (
	cc     = cache.Get()
	bus    = event.Get(false)
	logger = log.Get()
)
