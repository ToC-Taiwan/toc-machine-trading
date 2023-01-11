// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	dict map[string]*cache.Cache
	lock sync.RWMutex
}

func New() *Cache {
	return &Cache{
		dict: make(map[string]*cache.Cache),
	}
}

func (c *Cache) Set(k *Key, x interface{}) {
	if cc := c.getCacher(k); cc != nil {
		cc.Set(k.String(), x, 0)
		return
	}

	cc := cache.New(0, 0)
	cc.Set(k.String(), x, 0)
	c.addCacher(k, cc)
}

func (c *Cache) Get(k *Key) (interface{}, bool) {
	if cc := c.getCacher(k); cc != nil {
		return cc.Get(k.String())
	}
	return nil, false
}

func (c *Cache) getCacher(key *Key) *cache.Cache {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.dict[key.category]
}

func (c *Cache) addCacher(key *Key, cc *cache.Cache) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dict[key.category] = cc
}
