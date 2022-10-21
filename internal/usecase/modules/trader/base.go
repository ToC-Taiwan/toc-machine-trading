// Package trader package trader
package trader

import (
	"tmt/internal/usecase/modules/cache"
	"tmt/pkg/eventbus"
)

var (
	cc  = cache.GetCache()
	bus = eventbus.New()
)
