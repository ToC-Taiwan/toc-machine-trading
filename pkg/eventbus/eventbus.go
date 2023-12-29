// Package eventbus package eventbus
package eventbus

import (
	"sync"

	"github.com/asaskevich/EventBus"
	"github.com/patrickmn/go-cache"
)

var (
	singleton *Bus
	once      sync.Once
)

// Bus Bus
type Bus struct {
	bus EventBus.Bus
	cc  *cache.Cache
}

// Get if route is empty, return singleton
func Get(route ...string) *Bus {
	if singleton == nil {
		once.Do(func() {
			singleton = &Bus{
				bus: EventBus.New(),
				cc:  cache.New(0, 0),
			}
		})
		return Get(route...)
	}

	if len(route) == 0 {
		return singleton
	}

	bus := singleton
	for i := 0; i < len(route); i++ {
		if bus.getRoute(route[i]) == nil {
			v := &Bus{
				bus: EventBus.New(),
				cc:  cache.New(0, 0),
			}
			bus.addRoute(route[i], v)
		}
		bus = bus.getRoute(route[i])
	}
	return bus
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

func (c *Bus) addRoute(key string, value *Bus) {
	if _, ok := c.cc.Get(key); ok {
		return
	}
	c.cc.Set(key, value, cache.NoExpiration)
}

func (c *Bus) getRoute(key string) *Bus {
	if v, ok := c.cc.Get(key); ok {
		return v.(*Bus)
	}
	return nil
}
