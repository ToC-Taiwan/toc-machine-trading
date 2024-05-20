// Package mqtt package mqtt
package mqtt

import (
	"context"

	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-trade-protobuf/golang/pb"
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
