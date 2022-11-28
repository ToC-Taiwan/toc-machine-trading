// Package rabbit package rabbit
package rabbit

import (
	"fmt"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/modules/config"
	"tmt/internal/usecase/modules/logger"
	"tmt/pb"
	"tmt/pkg/common"
	"tmt/pkg/rabbitmq"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

var log = logger.Get()

const (
	routingKeyEvent        = "event"
	routingKeyOrder        = "order"
	routingKeyTick         = "tick"
	routingKeyFutureTick   = "future_tick"
	routingKeyBidAsk       = "bid_ask"
	routingKeyFutureBidAsk = "future_bid_ask"
)

// StreamRabbit -.
type StreamRabbit struct {
	conn *rabbitmq.Connection

	allStockMap  map[string]*entity.Stock
	allFutureMap map[string]*entity.Future

	futureTickChan map[string]chan *entity.RealTimeFutureTick
	futureTickLock sync.RWMutex

	orderStatusChanMap     map[string]chan interface{}
	orderStatusChanMapLock sync.RWMutex
}

// NewStream -.
func NewStream() *StreamRabbit {
	allConfig := config.GetConfig()

	conn := rabbitmq.NewConnection(
		allConfig.RabbitMQ.Exchange,
		allConfig.RabbitMQ.URL,
		allConfig.RabbitMQ.WaitTime,
		allConfig.RabbitMQ.Attempts,
	)

	if err := conn.AttemptConnect(); err != nil {
		log.Error(err)
	}

	return &StreamRabbit{
		conn:               conn,
		futureTickChan:     make(map[string]chan *entity.RealTimeFutureTick),
		orderStatusChanMap: make(map[string]chan interface{}),
	}
}

func (c *StreamRabbit) establishDelivery(key string) <-chan amqp.Delivery {
	delivery, err := c.conn.BindAndConsume(key)
	if err != nil {
		log.Panic(err)
	}
	return delivery
}

// FillAllBasic -.
func (c *StreamRabbit) FillAllBasic(allStockMap map[string]*entity.Stock, allFutureMap map[string]*entity.Future) {
	c.allStockMap = allStockMap
	c.allFutureMap = allFutureMap
}

// EventConsumer -.
func (c *StreamRabbit) EventConsumer(eventChan chan *entity.SinopacEvent) {
	delivery := c.establishDelivery(routingKeyEvent)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("EventConsumer rabbitMQ is closed")
			return
		}

		body := pb.EventMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetEventTime(), time.Local)
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

func (c *StreamRabbit) AddOrderStatusChan(orderStatusChan chan interface{}, connectionID string) {
	defer c.orderStatusChanMapLock.Unlock()
	c.orderStatusChanMapLock.Lock()
	c.orderStatusChanMap[connectionID] = orderStatusChan
}

func (c *StreamRabbit) RemoveOrderStatusChan(connectionID string) {
	defer c.orderStatusChanMapLock.Unlock()
	c.orderStatusChanMapLock.Lock()
	close(c.orderStatusChanMap[connectionID])
	delete(c.orderStatusChanMap, connectionID)
}

// OrderStatusConsumer OrderStatusConsumer
func (c *StreamRabbit) OrderStatusConsumer(orderStatusChan chan interface{}) {
	if len(c.allStockMap) == 0 || len(c.allFutureMap) == 0 {
		log.Panic("allStockMap or allFutureMap is empty")
	}
	c.orderStatusChanMapLock.Lock()
	c.orderStatusChanMap[uuid.New().String()] = orderStatusChan
	c.orderStatusChanMapLock.Unlock()
	delivery := c.establishDelivery(routingKeyOrder)
	for {
		d, opened := <-delivery
		if !opened {
			log.Error("OrderStatusConsumer rabbitMQ is closed")
			return
		}
		body := pb.StockOrderStatus{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}
		orderTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetOrderTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		var order interface{}
		switch {
		case c.allStockMap[body.GetCode()] != nil:
			order = &entity.StockOrder{
				StockNum: body.GetCode(),
				BaseOrder: entity.BaseOrder{
					OrderID:   body.GetOrderId(),
					Action:    entity.StringToOrderAction(body.GetAction()),
					Price:     body.GetPrice(),
					Quantity:  body.GetQuantity(),
					Status:    entity.StringToOrderStatus(body.GetStatus()),
					OrderTime: orderTime,
				},
			}
		case c.allFutureMap[body.GetCode()] != nil:
			order = &entity.FutureOrder{
				Code: body.GetCode(),
				BaseOrder: entity.BaseOrder{
					OrderID:   body.GetOrderId(),
					Action:    entity.StringToOrderAction(body.GetAction()),
					Price:     body.GetPrice(),
					Quantity:  body.GetQuantity(),
					Status:    entity.StringToOrderStatus(body.GetStatus()),
					OrderTime: orderTime,
				},
			}
		default:
			log.Warnf("unknown code: %s", body.GetCode())
		}

		c.orderStatusChanMapLock.RLock()
		for _, orderStatusChan := range c.orderStatusChanMap {
			orderStatusChan <- order
		}
		c.orderStatusChanMapLock.RUnlock()
	}
}

// TickConsumer -.
func (c *StreamRabbit) TickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyTick, stockNum))
	for {
		d, opened := <-delivery
		if !opened {
			log.Errorf("TickConsumer:%s rabbitMQ is closed", stockNum)
			return
		}

		body := pb.StockRealTimeTickMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() {
			continue
		}

		tickChan <- &entity.RealTimeStockTick{
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
		}
	}
}

// AddFutureTickChan -.
func (c *StreamRabbit) AddFutureTickChan(tickChan chan *entity.RealTimeFutureTick, connectionID string) {
	defer c.futureTickLock.Unlock()
	c.futureTickLock.Lock()
	c.futureTickChan[connectionID] = tickChan
}

// RemoveFutureTickChan -.
func (c *StreamRabbit) RemoveFutureTickChan(connectionID string) {
	defer c.futureTickLock.Unlock()
	c.futureTickLock.Lock()
	close(c.futureTickChan[connectionID])
	delete(c.futureTickChan, connectionID)
}

// FutureTickConsumer -.
func (c *StreamRabbit) FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick) {
	c.futureTickLock.Lock()
	c.futureTickChan[uuid.New().String()] = tickChan
	c.futureTickLock.Unlock()
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyFutureTick, code))
	for {
		d, opened := <-delivery
		if !opened {
			log.Errorf("FutureTickConsumer:%s rabbitMQ is closed", code)
			return
		}

		body := pb.FutureRealTimeTickMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() {
			continue
		}

		tick := &entity.RealTimeFutureTick{
			Code:            body.GetCode(),
			TickTime:        dataTime,
			Open:            body.GetOpen(),
			UnderlyingPrice: body.GetUnderlyingPrice(),
			BidSideTotalVol: body.GetBidSideTotalVol(),
			AskSideTotalVol: body.GetAskSideTotalVol(),
			AvgPrice:        body.GetAvgPrice(),
			Close:           body.GetClose(),
			High:            body.GetHigh(),
			Low:             body.GetLow(),
			Amount:          body.GetAmount(),
			TotalAmount:     body.GetTotalAmount(),
			Volume:          body.GetVolume(),
			TotalVolume:     body.GetTotalVolume(),
			TickType:        body.GetTickType(),
			ChgType:         body.GetChgType(),
			PriceChg:        body.GetPriceChg(),
			PctChg:          body.GetPctChg(),
		}

		c.futureTickLock.RLock()
		for _, tickChan := range c.futureTickChan {
			tickChan <- tick
		}
		c.futureTickLock.RUnlock()
	}
}

// StockBidAskConsumer -.
func (c *StreamRabbit) StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyBidAsk, stockNum))
	for {
		d, opened := <-delivery
		if !opened {
			log.Errorf("BidAskConsumer:%s rabbitMQ is closed", stockNum)
			return
		}

		body := pb.StockRealTimeBidAskMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() {
			continue
		}

		bidAskChan <- &entity.RealTimeStockBidAsk{
			StockNum: body.GetCode(),
			BidAskBase: entity.BidAskBase{
				BidAskTime: dataTime,
				BidPrice1:  body.GetBidPrice()[0], BidVolume1: body.GetBidVolume()[0], DiffBidVol1: body.GetDiffBidVol()[0],
				BidPrice2: body.GetBidPrice()[1], BidVolume2: body.GetBidVolume()[1], DiffBidVol2: body.GetDiffBidVol()[1],
				BidPrice3: body.GetBidPrice()[2], BidVolume3: body.GetBidVolume()[2], DiffBidVol3: body.GetDiffBidVol()[2],
				BidPrice4: body.GetBidPrice()[3], BidVolume4: body.GetBidVolume()[3], DiffBidVol4: body.GetDiffBidVol()[3],
				BidPrice5: body.GetBidPrice()[4], BidVolume5: body.GetBidVolume()[4], DiffBidVol5: body.GetDiffBidVol()[4],
				AskPrice1: body.GetAskPrice()[0], AskVolume1: body.GetAskVolume()[0], DiffAskVol1: body.GetDiffAskVol()[0],
				AskPrice2: body.GetAskPrice()[1], AskVolume2: body.GetAskVolume()[1], DiffAskVol2: body.GetDiffAskVol()[1],
				AskPrice3: body.GetAskPrice()[2], AskVolume3: body.GetAskVolume()[2], DiffAskVol3: body.GetDiffAskVol()[2],
				AskPrice4: body.GetAskPrice()[3], AskVolume4: body.GetAskVolume()[3], DiffAskVol4: body.GetDiffAskVol()[3],
				AskPrice5: body.GetAskPrice()[4], AskVolume5: body.GetAskVolume()[4], DiffAskVol5: body.GetDiffAskVol()[4],
			},
		}
	}
}

// FutureBidAskConsumer -.
func (c *StreamRabbit) FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyFutureBidAsk, code))
	for {
		d, opened := <-delivery
		if !opened {
			log.Errorf("FutureBidAskConsumer:%s rabbitMQ is closed", code)
			return
		}

		body := pb.FutureRealTimeBidAskMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			log.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(common.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			log.Error(err)
			continue
		}

		if body.GetSimtrade() {
			continue
		}

		bidAskChan <- &entity.FutureRealTimeBidAsk{
			Code: body.GetCode(),
			BidAskBase: entity.BidAskBase{
				BidAskTime: dataTime,
				BidPrice1:  body.GetBidPrice()[0], BidVolume1: body.GetBidVolume()[0], DiffBidVol1: body.GetDiffBidVol()[0],
				BidPrice2: body.GetBidPrice()[1], BidVolume2: body.GetBidVolume()[1], DiffBidVol2: body.GetDiffBidVol()[1],
				BidPrice3: body.GetBidPrice()[2], BidVolume3: body.GetBidVolume()[2], DiffBidVol3: body.GetDiffBidVol()[2],
				BidPrice4: body.GetBidPrice()[3], BidVolume4: body.GetBidVolume()[3], DiffBidVol4: body.GetDiffBidVol()[3],
				BidPrice5: body.GetBidPrice()[4], BidVolume5: body.GetBidVolume()[4], DiffBidVol5: body.GetDiffBidVol()[4],
				AskPrice1: body.GetAskPrice()[0], AskVolume1: body.GetAskVolume()[0], DiffAskVol1: body.GetDiffAskVol()[0],
				AskPrice2: body.GetAskPrice()[1], AskVolume2: body.GetAskVolume()[1], DiffAskVol2: body.GetDiffAskVol()[1],
				AskPrice3: body.GetAskPrice()[2], AskVolume3: body.GetAskVolume()[2], DiffAskVol3: body.GetDiffAskVol()[2],
				AskPrice4: body.GetAskPrice()[3], AskVolume4: body.GetAskVolume()[3], DiffAskVol4: body.GetDiffAskVol()[3],
				AskPrice5: body.GetAskPrice()[4], AskVolume5: body.GetAskVolume()[4], DiffAskVol5: body.GetDiffAskVol()[4],
			},
			BidTotalVol:          body.GetBidTotalVol(),
			AskTotalVol:          body.GetAskTotalVol(),
			UnderlyingPrice:      body.GetUnderlyingPrice(),
			FirstDerivedBidPrice: body.GetFirstDerivedBidPrice(),
			FirstDerivedAskPrice: body.GetFirstDerivedAskPrice(),
			FirstDerivedBidVol:   body.GetFirstDerivedBidVol(),
			FirstDerivedAskVol:   body.GetFirstDerivedAskVol(),
		}
	}
}
