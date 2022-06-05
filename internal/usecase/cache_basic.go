package usecase

import (
	"fmt"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

const (
	stockDetailKeyType string = "stock_detail"
)

// SetStockDetail SetStockDetail
func SetStockDetail(stock *entity.Stock) {
	key := cache.Key{
		Type: stockDetailKeyType,
		Name: fmt.Sprintf("%s:%s", stockDetailKeyType, stock.Number),
	}
	cache.Set(key, stock)
}

// GetStockDetail GetStockDetail
func GetStockDetail(stockNum string) *entity.Stock {
	key := cache.Key{
		Type: stockDetailKeyType,
		Name: fmt.Sprintf("%s:%s", stockDetailKeyType, stockNum),
	}

	if value, ok := cache.Get(key); ok {
		return value.(*entity.Stock)
	}

	return nil
}
