package usecase

import (
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

const (
	cacheCatagoryBasic cache.Category = "basic"
)

const (
	cacheIDStockDetail string = "stock_detail"
	cacheIDCalendar    string = "calendar"
	cacheIDBasicInfo   string = "basic_info"
)

var cc = New()

// GlobalCache -.
type GlobalCache struct {
	*cache.Cache
}

// New -.
func New() *GlobalCache {
	return &GlobalCache{
		Cache: cache.New(),
	}
}

// SetStockDetail SetStockDetail
func (c *GlobalCache) SetStockDetail(stock *entity.Stock) {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       fmt.Sprintf("%s:%s", cacheIDStockDetail, stock.Number),
	}
	c.Set(key, stock)
}

// GetStockDetail GetStockDetail
func (c *GlobalCache) GetStockDetail(stockNum string) *entity.Stock {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       fmt.Sprintf("%s:%s", cacheIDStockDetail, stockNum),
	}
	if value, ok := c.Get(key); ok {
		return value.(*entity.Stock)
	}
	return nil
}

// SetCalendar SetCalendar
func (c *GlobalCache) SetCalendar(calendar map[time.Time]bool) {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDCalendar,
	}
	c.Set(key, calendar)
}

// GetCalendar GetCalendar
func (c *GlobalCache) GetCalendar() map[time.Time]bool {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDCalendar,
	}
	if value, ok := c.Get(key); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}

// SetBasicInfo SetBasicInfo
func (c *GlobalCache) SetBasicInfo(info *entity.BasicInfo) {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDBasicInfo,
	}
	c.Set(key, info)
}

// GetBasicInfo GetBasicInfo
func (c *GlobalCache) GetBasicInfo() *entity.BasicInfo {
	key := cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDBasicInfo,
	}
	if value, ok := c.Get(key); ok {
		return value.(*entity.BasicInfo)
	}
	return nil
}
