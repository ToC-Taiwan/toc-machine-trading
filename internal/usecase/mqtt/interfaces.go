// Package mqtt package mqtt
package mqtt

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
)

type MQTT interface {
	EventConsumer(eventChan chan *entity.SinopacEvent)
	OrderStatusArrConsumer(orderStatusChan chan interface{})
	StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte)
	StockTickOddsPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte)
	FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick)
	FutureTickPbConsumer(ctx context.Context, code string, tickChan chan *pb.FutureRealTimeTickMessage)
	Unsubscribe(id int)
	Close()
}
