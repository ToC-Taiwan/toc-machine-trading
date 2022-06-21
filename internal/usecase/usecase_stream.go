package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"

	"github.com/google/go-cmp/cmp"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo   StreamRepo
	rabbit StreamRabbit
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *rabbit.StreamRabbit) {
	uc := &StreamUseCase{
		repo:   r,
		rabbit: t,
	}

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	if err := bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData); err != nil {
		log.Panic(err)
	}
}

// ReceiveEvent -.
func (uc *StreamUseCase) ReceiveEvent(ctx context.Context) {
	eventChan := make(chan *entity.SinopacEvent)
	go func() {
		for {
			event := <-eventChan
			if err := uc.repo.InsertEvent(ctx, event); err != nil {
				log.Error(err)
			}
		}
	}()
	uc.rabbit.EventConsumer(eventChan)
}

// ReceiveOrderStatus -.
func (uc *StreamUseCase) ReceiveOrderStatus(ctx context.Context) {
	orderStatusChan := make(chan *entity.Order)
	go func() {
		for {
			order := <-orderStatusChan
			cacheOrder := cc.GetOrderByOrderID(order.OrderID)
			if !cmp.Equal(order, cacheOrder) {
				cc.SetOrderByOrderID(order)
			}
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// ReceiveStreamData -.
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	for _, t := range targetArr {
		target := t
		tickChan := make(chan *entity.RealTimeTick)
		bidAskChan := make(chan *entity.RealTimeBidAsk)
		finishChan := make(chan struct{})
		go uc.tradeAgent(tickChan, bidAskChan, finishChan)
		for {
			_, ok := <-finishChan
			if !ok {
				break
			}
		}
		go uc.rabbit.TickConsumer(fmt.Sprintf("tick:%s", target.StockNum), tickChan)
		go uc.rabbit.BidAskConsumer(fmt.Sprintf("bid_ask:%s", target.StockNum), bidAskChan)
	}
	bus.PublishTopicEvent(topicSubscribeTickTargets, ctx, targetArr)
}

func (uc *StreamUseCase) tradeAgent(tickChan chan *entity.RealTimeTick, bidAskChan chan *entity.RealTimeBidAsk, finishChan chan struct{}) {
	var bidAsk *entity.RealTimeBidAsk
	var currentStatus orderMap
	go func() {
		for {
			tick := <-tickChan
			action := currentStatus.checkNeededPost()
			switch action {
			case entity.ActionNone:
				if bidAsk != nil && bidAsk.AskPrice1 == tick.Close {
					order := &entity.Order{
						StockNum:  tick.StockNum,
						Action:    entity.ActionBuy,
						Price:     tick.Close,
						Quantity:  1,
						TradeTime: time.Now(),
					}
					bus.PublishTopicEvent(topicOrder, order)
					go currentStatus.checkOrderStatus(order)
				}
			case entity.ActionWait:
				continue
			}
		}
	}()
	go func() {
		for {
			bidAsk = <-bidAskChan
		}
	}()

	// close channel to start receive data from rabbitmq
	close(finishChan)
}

type orderMap struct {
	data     map[entity.OrderAction][]*entity.Order
	mu       sync.RWMutex
	checking bool
}

func (o *orderMap) checkOrderStatus(order *entity.Order) {
	o.checking = true
	for {
		if order.OrderID != "" {
			break
		}
		time.Sleep(time.Second)
	}
	o.mu.Lock()
	o.data[order.Action] = append(o.data[order.Action], order)
	o.mu.Unlock()
	o.checking = false
}

func (o *orderMap) checkNeededPost() entity.OrderAction {
	if o.checking {
		return entity.ActionWait
	}
	o.mu.RLock()
	defer o.mu.RUnlock()
	if len(o.data[entity.ActionBuy]) > len(o.data[entity.ActionSell]) {
		return entity.ActionSell
	}

	if len(o.data[entity.ActionSellFirst]) > len(o.data[entity.ActionBuyLater]) {
		return entity.ActionBuyLater
	}
	return entity.ActionNone
}
