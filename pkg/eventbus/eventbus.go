// Package eventbus package eventbus
package eventbus

import (
	"github.com/asaskevich/EventBus"
)

// Bus Bus
type Bus struct {
	bus EventBus.Bus
}

// New New
func New() *Bus {
	return &Bus{
		bus: EventBus.New(),
	}
}
