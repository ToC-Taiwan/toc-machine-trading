// Package usecase package usecase
package usecase

import (
	"tmt/internal/entity"
)

//go:generate mockgen -source=interfaces_rabbit.go -destination=./mocks_rabbit_test.go -package=usecase_test

type Rabbit interface {
	FillAllBasic(allStockMap map[string]*entity.Stock, allFutureMap map[string]*entity.Future)

	EventConsumer(eventChan chan *entity.SinopacEvent)
	OrderStatusConsumer()
	OrderStatusArrConsumer()

	StockTickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick)
	StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk)

	FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk)
	AddFutureTickChan(tickChan chan *entity.RealTimeFutureTick, connectionID string)
	RemoveFutureTickChan(connectionID string)
	AddOrderStatusChan(orderStatusChan chan interface{}, connectionID string)
	RemoveOrderStatusChan(connectionID string)
}
