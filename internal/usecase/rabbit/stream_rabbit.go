// Package rabbit package rabbit
package rabbit

import (
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/rabbitmq"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

var log = logger.Get()

// StreamRabbit -.
type StreamRabbit struct {
	conn *rabbitmq.Connection
}

// NewStream -.
func NewStream() *StreamRabbit {
	allConfig, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	conn := rabbitmq.NewConnection(
		allConfig.RabbitMQ.Exchange,
		allConfig.RabbitMQ.URL,
		allConfig.RabbitMQ.WaitTime,
		allConfig.RabbitMQ.Attempts,
	)

	if err := conn.AttemptConnect(); err != nil {
		log.Error(err)
	}

	return &StreamRabbit{conn}
}

func (c *StreamRabbit) establishDelivery(key string) <-chan amqp.Delivery {
	delivery, err := c.conn.BindAndConsume(key)
	if err != nil {
		log.Panic(err)
	}
	return delivery
}

// EventConsumer -.
func (c *StreamRabbit) EventConsumer(eventChan chan *entity.SinopacEvent) {
	key := "event"
	delivery := c.establishDelivery(key)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("RabbitMQ is closed")
			return
		}

		body := pb.EventResponse{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(global.LongTimeLayout, body.GetEventTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		eventChan <- &entity.SinopacEvent{
			Event:     body.GetEvent(),
			EventCode: body.GetEventCode(),
			Info:      body.GetInfo(),
			Response:  body.GetRespCode(),
			EventTime: dataTime,
		}
	}
}

// OrderStatusConsumer OrderStatusConsumer
func (c *StreamRabbit) OrderStatusConsumer(orderStatusChan chan *entity.OrderStatus) {
	key := "order_status"
	delivery := c.establishDelivery(key)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("RabbitMQ is closed")
			return
		}

		body := pb.StockOrderStatus{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		actionMap := entity.ActionListMap
		statusMap := entity.StatusListMap
		orderTime, err := time.ParseInLocation(global.LongTimeLayout, body.GetOrderTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}
		orderStatusChan <- &entity.OrderStatus{
			StockNum:  body.GetCode(),
			Action:    actionMap[body.GetAction()],
			Price:     body.GetPrice(),
			Quantity:  body.GetQuantity(),
			Status:    statusMap[body.GetStatus()],
			OrderID:   body.GetOrderId(),
			OrderTime: orderTime,
		}
	}
}

// TickConsumer -.
func (c *StreamRabbit) TickConsumer(key string, tickChan chan *entity.RealTimeTick) {
	delivery := c.establishDelivery(key)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("RabbitMQ is closed")
			return
		}

		body := pb.StockRealTimeTickResponse{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(global.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() == 1 {
			continue
		}

		tickChan <- &entity.RealTimeTick{
			StockNum:        body.GetCode(),
			TickTime:        dataTime,
			Open:            body.GetOpen(),
			AvgPrice:        body.GetAvgPrice(),
			Close:           body.GetClose(),
			High:            body.GetHigh(),
			Low:             body.GetLow(),
			Amount:          body.GetAmount(),
			AmountSum:       body.GetTotalAmount(),
			Volume:          body.GetVolume(),
			VolumeSum:       body.GetTotalVolume(),
			TickType:        body.GetTickType(),
			ChgType:         body.GetChgType(),
			PriceChg:        body.GetPriceChg(),
			PctChg:          body.GetPctChg(),
			BidSideTotalVol: body.GetBidSideTotalVol(),
			AskSideTotalVol: body.GetAskSideTotalVol(),
			BidSideTotalCnt: body.GetBidSideTotalCnt(),
			AskSideTotalCnt: body.GetAskSideTotalCnt(),
			Suspend:         body.GetSuspend(),
			Simtrade:        body.GetSimtrade(),
		}
	}
}

// BidAskConsumer -.
func (c *StreamRabbit) BidAskConsumer(key string, bidAskChan chan *entity.RealTimeBidAsk) {
	delivery := c.establishDelivery(key)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("RabbitMQ is closed")
			return
		}

		body := pb.StockRealTimeBidAskResponse{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(global.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() == 1 {
			continue
		}

		bidAskChan <- &entity.RealTimeBidAsk{
			StockNum:  body.GetCode(),
			TickTime:  dataTime,
			BidPrice1: body.GetBidPrice()[0], BidVolume1: body.GetBidVolume()[0], DiffBidVol1: body.GetDiffBidVol()[0],
			BidPrice2: body.GetBidPrice()[1], BidVolume2: body.GetBidVolume()[1], DiffBidVol2: body.GetDiffBidVol()[1],
			BidPrice3: body.GetBidPrice()[2], BidVolume3: body.GetBidVolume()[2], DiffBidVol3: body.GetDiffBidVol()[2],
			BidPrice4: body.GetBidPrice()[3], BidVolume4: body.GetBidVolume()[3], DiffBidVol4: body.GetDiffBidVol()[3],
			BidPrice5: body.GetBidPrice()[4], BidVolume5: body.GetBidVolume()[4], DiffBidVol5: body.GetDiffBidVol()[4],
			AskPrice1: body.GetAskPrice()[0], AskVolume1: body.GetAskVolume()[0], DiffAskVol1: body.GetDiffAskVol()[0],
			AskPrice2: body.GetAskPrice()[1], AskVolume2: body.GetAskVolume()[1], DiffAskVol2: body.GetDiffAskVol()[1],
			AskPrice3: body.GetAskPrice()[2], AskVolume3: body.GetAskVolume()[2], DiffAskVol3: body.GetDiffAskVol()[2],
			AskPrice4: body.GetAskPrice()[3], AskVolume4: body.GetAskVolume()[3], DiffAskVol4: body.GetDiffAskVol()[3],
			AskPrice5: body.GetAskPrice()[4], AskVolume5: body.GetAskVolume()[4], DiffAskVol5: body.GetDiffAskVol()[4],
			Suspend:  body.GetSuspend(),
			Simtrade: body.GetSimtrade(),
		}
	}
}
