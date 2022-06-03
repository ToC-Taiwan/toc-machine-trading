// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

type cacheType string

// const (
// 	noExpired time.Duration = 0
// 	noCleanUp time.Duration = 0
// )

// Cache Cache
type Cache struct {
	CacheMap map[string]*cache.Cache
	lock     sync.RWMutex
}

// Key Key
type Key struct {
	Name string
	Type cacheType
}

var (
	globalCache *Cache
	once        sync.Once
)

// GetCache GetCache
func GetCache() *Cache {
	if globalCache != nil {
		return globalCache
	}
	once.Do(initGlobalCache)
	return globalCache
}

func initGlobalCache() {
	if globalCache != nil {
		return
	}
	var newCache Cache
	newCache.CacheMap = make(map[string]*cache.Cache)
	globalCache = &newCache
}

// getCacheByType getCacheByType
// func (c *Cache) getCacheByType(cacheType cacheType) *cache.Cache {
// 	c.lock.RLock()
// 	tmp := c.CacheMap[string(cacheType)]
// 	c.lock.RUnlock()
// 	if tmp == nil {
// 		tmp = cache.New(noExpired, noCleanUp)
// 		c.lock.Lock()
// 		c.CacheMap[string(cacheType)] = tmp
// 		c.lock.Unlock()
// 	}
// 	return tmp
// }

// GetAllCacheType GetAllCacheType
func (c *Cache) GetAllCacheType() []string {
	c.lock.RLock()
	var typeArr []string
	for k := range c.CacheMap {
		typeArr = append(typeArr, k)
	}
	c.lock.RUnlock()
	return typeArr
}

// GetAllCacheByType GetAllCacheByType
func (c *Cache) GetAllCacheByType(cacheType string) interface{} {
	defer c.lock.RUnlock()
	c.lock.RLock()
	return c.CacheMap[cacheType].Items()
}
