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

// StockTickConsumer -.
func (c *Rabbit) StockTickConsumer(stockNum string, tickChan chan *entity.RealTimeStockTick) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyTick, stockNum))
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := pb.StockRealTimeTickMessage{}
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

// StockBidAskConsumer -.
func (c *Rabbit) StockBidAskConsumer(stockNum string, bidAskChan chan *entity.RealTimeStockBidAsk) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyBidAsk, stockNum))
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := pb.StockRealTimeBidAskMessage{}
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
func (c *Rabbit) FutureBidAskConsumer(code string, bidAskChan chan *entity.FutureRealTimeBidAsk) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyFutureBidAsk, code))
	for {
		d, opened := <-delivery
		if !opened {
			return
		}

		body := pb.FutureRealTimeBidAskMessage{}
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

// StockTickPbConsumer -.
func (c *Rabbit) StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte) {
	delivery := c.establishDelivery(fmt.Sprintf("%s:%s", routingKeyTick, stockNum))
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
