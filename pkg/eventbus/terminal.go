package eventbus

import (
	"github.com/patrickmn/go-cache"
)

type busTerminal struct {
	terminal *cache.Cache
}

func newBusTerminal() *busTerminal {
	return &busTerminal{
		terminal: cache.New(0, 0),
	}
}

func (bt *busTerminal) addBus(key string, value *Bus) {
	bt.terminal.Set(key, value, cache.NoExpiration)
}

func (bt *busTerminal) getBus(key string) *Bus {
	if v, ok := bt.terminal.Get(key); ok {
		return v.(*Bus)
	}
	return nil
}
