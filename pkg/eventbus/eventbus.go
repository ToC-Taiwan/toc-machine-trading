// Package eventbus package eventbus
package eventbus

import (
	"toc-machine-trading/pkg/logger"

	"github.com/asaskevich/EventBus"
)

var log = logger.Get()

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

// PublishTopicEvent PublishTopicEvent
func (c *Bus) PublishTopicEvent(topic string, arg ...interface{}) {
	go c.bus.Publish(topic, arg...)
}

// SubscribeTopic SubscribeTopic
func (c *Bus) SubscribeTopic(topic string, fn ...interface{}) {
	for _, f := range fn {
		err := c.bus.SubscribeAsync(topic, f, false)
		if err != nil {
			log.Panic(err)
		}
	}
}
