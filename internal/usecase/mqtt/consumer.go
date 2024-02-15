package mqtt

import (
	"context"
	"fmt"
	"time"

	"tmt/internal/entity"
	"tmt/pb"

	"google.golang.org/protobuf/proto"
)

// EventConsumer -.
func (c *Rabbit) EventConsumer(eventChan chan *entity.SinopacEvent) {
	delivery := c.establishDelivery(routingKeyEvent)
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := pb.EventMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			c.logger.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(entity.LongTimeLayout, body.GetEventTime(), time.Local)
		if err != nil {
			c.logger.Error(err)
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

// OrderStatusConsumer _.
func (c *Rabbit) OrderStatusConsumer(orderStatusChan chan interface{}) {
	delivery := c.establishDelivery(routingKeyOrder)
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := &pb.OrderStatus{}
		if err := proto.Unmarshal(d.Body, body); err != nil {
			c.logger.Error(err)
			continue
		}

		if order := c.protoToOrder(body); order == nil {
			continue
		} else {
			orderStatusChan <- order
		}
	}
}

// OrderStatusArrConsumer -.
func (c *Rabbit) OrderStatusArrConsumer(orderStatusChan chan interface{}) {
	delivery := c.establishDelivery(routingKeyOrderArr)
	for {
		d, opened := <-delivery
		if !opened {
			return
		}
		body := &pb.OrderStatusArr{}
		if err := proto.Unmarshal(d.Body, body); err != nil {
			c.logger.Error(err)
			continue
		}

		for _, b := range body.GetData() {
			if data := c.protoToOrder(b); data != nil {
				orderStatusChan <- data
			}
		}
	}
}

// FutureTickConsumer -.
func (c *Rabbit) FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyFutureTick, code))
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := pb.FutureRealTimeTickMessage{}
		if err := proto.Unmarshal(d.Body, &body); err != nil {
			c.logger.Error(err)
			continue
		}

		dataTime, err := time.ParseInLocation(entity.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			c.logger.Error(err)
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

		tickChan <- tick
	}
}

// StockTickPbConsumer -.
func (c *Rabbit) StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyStockTick, stockNum))
	for {
		select {
		case <-ctx.Done():
			return

		case d, opened := <-delivery:
			if !opened {
				return
			}
			tickChan <- d.Body
		}
	}
}

// StockTickOddsPbConsumer -.
func (c *Rabbit) StockTickOddsPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyStockTickOdds, stockNum))
	for {
		select {
		case <-ctx.Done():
			return

		case d, opened := <-delivery:
			if !opened {
				return
			}
			tickChan <- d.Body
		}
	}
}
