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

// PublishTopicEvent PublishTopicEvent
func (c *Bus) PublishTopicEvent(topic string, arg ...interface{}) {
	go c.bus.Publish(topic, arg...)
}

// SubscribeTopic SubscribeTopic
func (c *Bus) SubscribeTopic(topic string, fn interface{}) {
	err := c.bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		panic(err)
	}
}

func (c *Bus) UnSubscribeTopic(topic string, fn interface{}) {
	err := c.bus.Unsubscribe(topic, fn)
	if err != nil {
		panic(err)
	}
}
