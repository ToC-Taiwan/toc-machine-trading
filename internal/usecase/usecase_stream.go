package usecase

import (
	"context"
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/rabbit"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo   StreamRepo
	rabbit StreamRabbit
	bus    *eventbus.Bus
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *rabbit.StreamRabbit, bus *eventbus.Bus) {
	uc := &StreamUseCase{
		repo:   r,
		rabbit: t,
		bus:    bus,
	}

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	if err := uc.bus.SubscribeTopic(topicStreamTargets, uc.ReceiveTicks); err != nil {
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
	orderStatusChan := make(chan *entity.OrderStatus)
	go func() {
		for {
			orderStatus := <-orderStatusChan
			log.Info(orderStatus)
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// ReceiveTicks -.
func (uc *StreamUseCase) ReceiveTicks(ctx context.Context, targetArr []*entity.Target) {
	for _, t := range targetArr {
		tickChan := make(chan *entity.RealTimeTick)
		bidAskChan := make(chan *entity.RealTimeBidAsk)
		go uc.tradeAgent(tickChan, bidAskChan)
		go uc.rabbit.TickConsumer(fmt.Sprintf("tick:%s", t.StockNum), tickChan)
		go uc.rabbit.BidAskConsumer(fmt.Sprintf("bid_ask:%s", t.StockNum), bidAskChan)
	}
	uc.bus.PublishTopicEvent(topicSubscribeTickTargets, ctx, targetArr)
}

func (uc *StreamUseCase) tradeAgent(tickChan chan *entity.RealTimeTick, bidAskChan chan *entity.RealTimeBidAsk) {
	wait := make(chan struct{})
	go func() {
		for {
			tick := <-tickChan
			log.Infof("tick:%s\n", time.Since(tick.TickTime).String())
		}
	}()
	go func() {
		for {
			bidAsk := <-bidAskChan
			log.Infof("bidask:%s\n", time.Since(bidAsk.TickTime).String())
		}
	}()
	<-wait
}
