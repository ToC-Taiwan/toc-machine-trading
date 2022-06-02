// Package eventbus package eventbus
package eventbus

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

// Bus Bus
type Bus struct {
	bus EventBus.Bus
}

var (
	globalBus *Bus
	once      sync.Once
)

func initBus() {
	if globalBus != nil {
		return
	}
	newAgent := &Bus{
		bus: EventBus.New(),
	}
	globalBus = newAgent
}

// Get Get
func Get() *Bus {
	if globalBus != nil {
		return globalBus
	}
	once.Do(initBus)
	return globalBus
}
