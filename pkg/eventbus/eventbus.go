// Package eventbus package eventbus
package eventbus

import (
	"github.com/asaskevich/EventBus"
)

var (
	singleton *Bus
	terminal  *busTerminal = newBusTerminal()
)

// Bus Bus
type Bus struct {
	bus EventBus.Bus
}

func Get(route ...string) *Bus {
	if singleton == nil {
		singleton = &Bus{
			bus: EventBus.New(),
		}
	}

	switch len(route) {
	case 0:
		return singleton
	case 1:
		if v := terminal.getBus(route[0]); v != nil {
			return v
		}

		bus := &Bus{
			bus: EventBus.New(),
		}
		terminal.addBus(route[0], bus)
		return bus
	default:
		panic("route length must be 0 or 1")
	}
}

func (c *Bus) PublishTopicEvent(topic string, arg ...interface{}) {
	c.bus.Publish(topic, arg...)
}

// SubscribeAsync Transactional determines whether subsequent callbacks for a topic are
// run serially (true) or concurrently (false)
func (c *Bus) SubscribeAsync(topic string, transactional bool, fn ...interface{}) {
	if len(fn) == 0 {
		panic("fn length must be greater than 0")
	}

	for i := 0; i < len(fn); i++ {
		err := c.bus.SubscribeAsync(topic, fn[i], transactional)
		if err != nil {
			panic(err)
		}
	}
}

// Subscribe -.
func (c *Bus) Subscribe(topic string, fn ...interface{}) {
	if len(fn) == 0 {
		panic("fn length must be greater than 0")
	}

	for i := 0; i < len(fn); i++ {
		err := c.bus.Subscribe(topic, fn[i])
		if err != nil {
			panic(err)
		}
	}
}

func (c *Bus) UnSubscribe(topic string, fn ...interface{}) {
	if len(fn) == 0 {
		panic("fn length must be greater than 0")
	}

	for i := len(fn) - 1; i >= 0; i-- {
		err := c.bus.Unsubscribe(topic, fn[i])
		if err != nil {
			panic(err)
		}
	}
}
