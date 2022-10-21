// Package eventbus package eventbus
package eventbus

import (
	"github.com/asaskevich/EventBus"
)

var singleInstance *Bus

// Bus Bus
type Bus struct {
	bus EventBus.Bus
}

// New New
func New() *Bus {
	if singleInstance != nil {
		return singleInstance
	}

	singleInstance = &Bus{
		bus: EventBus.New(),
	}

	return singleInstance
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
			panic(err)
		}
	}
}
