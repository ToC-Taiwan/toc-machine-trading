// Package cache package cache
package cache

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	dict map[int64]*cache.Cache
	lock sync.RWMutex
}

func New() *Cache {
	return &Cache{
		dict: make(map[int64]*cache.Cache),
	}
}

func (c *Cache) splitKey(k string) (int64, string) {
	before, after, ok := strings.Cut(k, ":")
	if !ok {
		panic(fmt.Sprintf("invalid key: %s", k))
	}

	category, err := strconv.ParseInt(before, 10, 8)
	if err != nil {
		panic(fmt.Sprintf("invalid category: %s", before))
	}
	return category, after
}

func (c *Cache) getCacher(category int64) *cache.Cache {
	c.lock.RLock()
	cc := c.dict[category]
	c.lock.RUnlock()
	if cc != nil {
		return cc
	}

	c.lock.Lock()
	c.dict[category] = cache.New(0, 0)
	c.lock.Unlock()
	return c.dict[category]
}

func (c *Cache) Set(k string, x interface{}) {
	category, k := c.splitKey(k)
	c.getCacher(category).Set(k, x, 0)
}

func (c *Cache) Get(k string) (interface{}, bool) {
	category, k := c.splitKey(k)
	return c.getCacher(category).Get(k)
}

func (c *Cache) GetAll(category int64) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range c.getCacher(category).Items() {
		result[k] = v.Object
	}
	return result
}
