package eventbus

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

type busTerminal struct {
	terminal *cache.Cache
	lock     sync.RWMutex
}

func newBusTerminal() *busTerminal {
	return &busTerminal{
		terminal: cache.New(0, 0),
	}
}

func (bt *busTerminal) addBus(key string, value *Bus) {
	defer bt.lock.Unlock()
	bt.lock.Lock()

	if _, ok := bt.terminal.Get(key); ok {
		return
	}

	bt.terminal.Set(key, value, cache.NoExpiration)
}

func (bt *busTerminal) getBus(key string) *Bus {
	defer bt.lock.RUnlock()
	bt.lock.RLock()

	if v, ok := bt.terminal.Get(key); ok {
		return v.(*Bus)
	}
	return nil
}
