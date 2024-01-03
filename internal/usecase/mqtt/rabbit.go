// Package mqtt package mqtt
package mqtt

import (
	"time"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/log"
	"tmt/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Rabbit -.
type Rabbit struct {
	conn   *rabbitmq.Connection
	logger *log.Log
}

func NewRabbit(connection *rabbitmq.Connection) *Rabbit {
	return &Rabbit{
		conn:   connection,
		logger: log.Get(),
	}
}

func (c *Rabbit) Close() {
	if e := c.conn.Close(); e != nil {
		c.logger.Error(e)
	}
}

func (c *Rabbit) establishDelivery(key string) <-chan amqp.Delivery {
	delivery, err := c.conn.BindAndConsume(key)
	if err != nil {
		c.logger.Fatal(err)
	}
	return delivery
}

func (c *Rabbit) protoToOrder(proto *pb.OrderStatus) interface{} {
	orderTime, err := time.ParseInLocation(entity.LongTimeLayout, proto.GetOrderTime(), time.Local)
	if err != nil {
		c.logger.Error(err)
		return nil
	}

	detail := entity.OrderDetail{
		OrderID:   proto.GetOrderId(),
		Action:    entity.StringToOrderAction(proto.GetAction()),
		Price:     proto.GetPrice(),
		Status:    entity.StringToOrderStatus(proto.GetStatus()),
		OrderTime: orderTime,
	}

	switch proto.GetType() {
	case pb.OrderType_TYPE_STOCK_LOT:
		return &entity.StockOrder{
			StockNum:    proto.GetCode(),
			Lot:         proto.GetQuantity(),
			OrderDetail: detail,
		}
	case pb.OrderType_TYPE_STOCK_SHARE:
		return &entity.StockOrder{
			StockNum:    proto.GetCode(),
			Share:       proto.GetQuantity(),
			OrderDetail: detail,
		}
	case pb.OrderType_TYPE_FUTURE:
		return &entity.FutureOrder{
			Code:        proto.GetCode(),
			Position:    proto.GetQuantity(),
			OrderDetail: detail,
		}
	default:
		c.logger.Warnf("protoToOrder: unknown order type")
		return nil
	}
}
