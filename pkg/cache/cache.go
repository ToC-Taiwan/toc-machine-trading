// Package cache package cache
package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	noExpired time.Duration = 0
	noCleanUp time.Duration = 0
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

var (
	globalCache *Cache
	once        sync.Once
)

func initGlobalCache() {
	if globalCache != nil {
		return
	}
	var newCache Cache
	newCache.CacheMap = make(map[string]*cache.Cache)
	globalCache = &newCache
}

func getCacheByType(keyType string) *cache.Cache {
	if globalCache == nil {
		once.Do(initGlobalCache)
	}

	globalCache.lock.RLock()
	tmp := globalCache.CacheMap[keyType]
	globalCache.lock.RUnlock()

	if tmp == nil {
		tmp = cache.New(noExpired, noCleanUp)
		globalCache.lock.Lock()
		globalCache.CacheMap[keyType] = tmp
		globalCache.lock.Unlock()
	}
	return tmp
}

// Set -.
func Set(k Key, x interface{}) {
	getCacheByType(k.Type).Set(k.Name, x, noExpired)
}

// Get -.
func Get(k Key) (interface{}, bool) {
	return getCacheByType(k.Type).Get(k.Name)
}
