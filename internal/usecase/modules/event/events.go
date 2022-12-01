package event

import "tmt/pkg/eventbus"

var singletonBus *Bus

// Bus -.
type Bus struct {
	bus *eventbus.Bus
}

// Get -.
func Get(new bool) *Bus {
	if singletonBus != nil && !new {
		return singletonBus
	}

	singletonBus = &Bus{
		bus: eventbus.New(),
	}

	return singletonBus
}

// PublishTopicEvent -.
func (b *Bus) PublishTopicEvent(topic string, arg ...interface{}) {
	b.bus.PublishTopicEvent(topic, arg...)
}

// SubscribeTopic -.
func (b *Bus) SubscribeTopic(topic string, fn ...interface{}) {
	for _, f := range fn {
		b.bus.SubscribeTopic(topic, f)
	}
}

// UnSubscribeTopic -.
func (b *Bus) UnSubscribeTopic(topic string, fn interface{}) {
	b.bus.UnSubscribeTopic(topic, fn)
}
