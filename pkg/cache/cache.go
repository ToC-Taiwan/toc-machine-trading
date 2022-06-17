// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

var (
	globalCache *Cache
	once        sync.Once
)

// Cache Cache
type Cache struct {
	CacheMap map[string]*cache.Cache
	lock     sync.RWMutex
}

// Key Key
type Key struct {
	Type string
	Name string
}

func initGlobalCache() {
	if globalCache != nil {
		return
	}
	var newCache Cache
	newCache.CacheMap = make(map[string]*cache.Cache)
	globalCache = &newCache
}

func getOrCreateCache(keyType string) *cache.Cache {
	if globalCache == nil {
		once.Do(initGlobalCache)
	}

	globalCache.lock.RLock()
	tmp := globalCache.CacheMap[keyType]
	globalCache.lock.RUnlock()

	if tmp == nil {
		tmp = cache.New(0, 0)
		globalCache.lock.Lock()
		globalCache.CacheMap[keyType] = tmp
		globalCache.lock.Unlock()
	}
	return tmp
}

// Set -.
func Set(k Key, x interface{}) {
	getOrCreateCache(k.Type).Set(k.Name, x, 0)
}

// Get -.
func Get(k Key) (interface{}, bool) {
	return getOrCreateCache(k.Type).Get(k.Name)
}
