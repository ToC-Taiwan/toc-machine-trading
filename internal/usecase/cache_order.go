package usecase

import (
	"fmt"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

const (
	keyTypeTickChan   string = "tick_chan"
	keyTypeBidAskChan string = "bidask_chan"
)

// CacheSetTickChan -.
func CacheSetTickChan(stockNum string, tickChan chan *entity.RealTimeTick) {
	key := cache.Key{
		Type: keyTypeTickChan,
		Name: fmt.Sprintf("%s:%s", keyTypeTickChan, stockNum),
	}
	cache.Set(key, tickChan)
}

// CacheGetTickChan -.
func CacheGetTickChan(stockNum string) chan *entity.RealTimeTick {
	key := cache.Key{
		Type: keyTypeTickChan,
		Name: fmt.Sprintf("%s:%s", keyTypeTickChan, stockNum),
	}
	if value, ok := cache.Get(key); ok {
		return value.(chan *entity.RealTimeTick)
	}
	return nil
}

// CacheSetBidAskChan -.
func CacheSetBidAskChan(stockNum string, bidAskChan chan *entity.RealTimeBidAsk) {
	key := cache.Key{
		Type: keyTypeBidAskChan,
		Name: fmt.Sprintf("%s:%s", keyTypeTickChan, stockNum),
	}
	cache.Set(key, bidAskChan)
}

// CacheGetBidAskChan -.
func CacheGetBidAskChan(stockNum string) chan *entity.RealTimeBidAsk {
	key := cache.Key{
		Type: keyTypeBidAskChan,
		Name: fmt.Sprintf("%s:%s", keyTypeTickChan, stockNum),
	}
	if value, ok := cache.Get(key); ok {
		return value.(chan *entity.RealTimeBidAsk)
	}
	fmt.Println("not ok")
	return nil
}
