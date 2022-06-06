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
func (c *Bus) PublishTopicEvent(topic string, event interface{}) {
	go c.bus.Publish(topic, event)
}

// SubscribeTopic SubscribeTopic
func (c *Bus) SubscribeTopic(topic string, f func(opt ...interface{})) error {
	err := c.bus.Subscribe(topic, f)
	if err != nil {
		return err
	}
	return nil
}
