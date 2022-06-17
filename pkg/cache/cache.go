// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

// Cache Cache
type Cache struct {
	CacheMap map[string]*cache.Cache
	lock     sync.RWMutex
}

// Category -.
type Category string

func (c Category) String() string {
	return string(c)
}

// Key Key
type Key struct {
	Category Category
	ID       string
}

// New -.
func New() *Cache {
	newCache := &Cache{}
	newCache.CacheMap = make(map[string]*cache.Cache)
	return newCache
}

func (c *Cache) getOrNewCache(category string) *cache.Cache {
	c.lock.RLock()
	cc := c.CacheMap[category]
	c.lock.RUnlock()

	if cc == nil {
		cc = cache.New(0, 0)
		c.lock.Lock()
		c.CacheMap[category] = cc
		c.lock.Unlock()
	}
	return cc
}

// Set -.
func (c *Cache) Set(k Key, x interface{}) {
	c.getOrNewCache(k.Category.String()).Set(k.ID, x, 0)
}

// Get -.
func (c *Cache) Get(k Key) (interface{}, bool) {
	return c.getOrNewCache(k.Category.String()).Get(k.ID)
}
