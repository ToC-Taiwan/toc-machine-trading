// Package trader package trader
package trader

import (
	"tmt/internal/usecase/cache"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

var (
	cc     = cache.Get()
	bus    = eventbus.Get()
	logger = log.Get()
)
