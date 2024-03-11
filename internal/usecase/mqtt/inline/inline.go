// Package inline package inline
package inline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/mqtt"
	"tmt/pb"
	"tmt/pkg/embedbkr"

	mqttSrv "github.com/mochi-mqtt/server/v2"

	"github.com/mochi-mqtt/server/v2/packets"
	"google.golang.org/protobuf/proto"
)

type Inliner struct {
	srv *embedbkr.MQSrv

	subIDMap     map[int]struct{}
	subIDMapLock sync.Mutex
}

func NewInliner() mqtt.MQTT {
	return &Inliner{
		srv:      embedbkr.Get(),
		subIDMap: make(map[int]struct{}),
	}
}

func (i *Inliner) addID(id int) {
	i.subIDMapLock.Lock()
	defer i.subIDMapLock.Unlock()
	i.subIDMap[id] = struct{}{}
}

func (i *Inliner) Close() {
	i.subIDMapLock.Lock()
	defer i.subIDMapLock.Unlock()
	for id := range i.subIDMap {
		i.srv.Unsubscribe(id)
	}
}

func (i *Inliner) Unsubscribe(id int) {
	i.subIDMapLock.Lock()
	defer i.subIDMapLock.Unlock()
	delete(i.subIDMap, id)
	i.srv.Unsubscribe(id)
}

func (i *Inliner) EventConsumer(eventChan chan *entity.SinopacEvent) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		body := pb.EventMessage{}
		if err := proto.Unmarshal(pk.Payload, &body); err != nil {
			return
		}

		dataTime, err := time.ParseInLocation(entity.LongTimeLayout, body.GetEventTime(), time.Local)
		if err != nil {
			return
		}

		eventChan <- &entity.SinopacEvent{
			Event:     body.GetEvent(),
			EventCode: body.GetEventCode(),
			Info:      body.GetInfo(),
			Response:  body.GetRespCode(),
			EventTime: dataTime,
		}
	}
	topic := fmt.Sprintf("direct/%s", mqtt.RoutingKeyEvent)
	if id := i.srv.Subscribe(topic, callbackFn); id != -1 {
		i.addID(id)
	}
}

func (i *Inliner) OrderStatusArrConsumer(orderStatusChan chan interface{}) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		body := &pb.OrderStatusArr{}
		if err := proto.Unmarshal(pk.Payload, body); err != nil {
			return
		}

		for _, b := range body.GetData() {
			if data := i.protoToOrder(b); data != nil {
				orderStatusChan <- data
			}
		}
	}
	topic := fmt.Sprintf("direct/%s", mqtt.RoutingKeyOrderArr)
	if id := i.srv.Subscribe(topic, callbackFn); id != -1 {
		i.addID(id)
	}
}

func (i *Inliner) StockTickPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		tickChan <- pk.Payload
	}
	topic := fmt.Sprintf("direct/%s/%s", mqtt.RoutingKeyStockTick, stockNum)
	id := i.srv.Subscribe(topic, callbackFn)
	if id != -1 {
		i.addID(id)
	}
	<-ctx.Done()
	i.Unsubscribe(id)
}

func (i *Inliner) StockTickOddsPbConsumer(ctx context.Context, stockNum string, tickChan chan []byte) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		tickChan <- pk.Payload
	}
	topic := fmt.Sprintf("direct/%s/%s", mqtt.RoutingKeyStockTickOdds, stockNum)
	id := i.srv.Subscribe(topic, callbackFn)
	if id != -1 {
		i.addID(id)
	}
	<-ctx.Done()
	i.Unsubscribe(id)
}

func (i *Inliner) FutureTickConsumer(code string, tickChan chan *entity.RealTimeFutureTick) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		body := pb.FutureRealTimeTickMessage{}
		if err := proto.Unmarshal(pk.Payload, &body); err != nil {
			return
		}

		dataTime, err := time.ParseInLocation(entity.LongTimeLayout, body.GetDateTime(), time.Local)
		if err != nil {
			return
		}

		if body.GetSimtrade() {
			return
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
	topic := fmt.Sprintf("direct/%s/%s", mqtt.RoutingKeyFutureTick, code)
	if id := i.srv.Subscribe(topic, callbackFn); id != -1 {
		i.addID(id)
	}
}

func (i *Inliner) FutureTickPbConsumer(ctx context.Context, code string, tickChan chan *pb.FutureRealTimeTickMessage) {
	callbackFn := func(cl *mqttSrv.Client, sub packets.Subscription, pk packets.Packet) {
		body := pb.FutureRealTimeTickMessage{}
		if err := proto.Unmarshal(pk.Payload, &body); err != nil {
			return
		}

		tickChan <- &body
	}
	topic := fmt.Sprintf("direct/%s/%s", mqtt.RoutingKeyFutureTick, code)
	id := i.srv.Subscribe(topic, callbackFn)
	if id != -1 {
		i.addID(id)
	}
	<-ctx.Done()
	i.Unsubscribe(id)
}

func (i *Inliner) protoToOrder(proto *pb.OrderStatus) interface{} {
	orderTime, err := time.ParseInLocation(entity.LongTimeLayout, proto.GetOrderTime(), time.Local)
	if err != nil {
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
		return nil
	}
}
