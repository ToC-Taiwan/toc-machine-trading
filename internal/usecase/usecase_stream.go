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

	if err := uc.bus.SubscribeTopic(topicStreamTickTargets, uc.ReceiveTicks); err != nil {
		log.Panic(err)
	}
	if err := uc.bus.SubscribeTopic(topicStreamBidAskTargets, uc.ReceiveBidAsk); err != nil {
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
		target := t
		go func() {
			tickChan := make(chan *entity.RealTimeTick)
			go tickProcessor(tickChan)

			uc.rabbit.TickConsumer(fmt.Sprintf("tick:%s", target.StockNum), tickChan)
		}()
	}
	uc.bus.PublishTopicEvent(topicSubscribeTickTargets, ctx, targetArr)
}

// ReceiveBidAsk -.
func (uc *StreamUseCase) ReceiveBidAsk(ctx context.Context, targetArr []*entity.Target) {
	for _, t := range targetArr {
		target := t
		go func() {
			bidAskChan := make(chan *entity.RealTimeBidAsk)
			go bidAskProcessor(bidAskChan)

			uc.rabbit.BidAskConsumer(fmt.Sprintf("bid_ask:%s", target.StockNum), bidAskChan)
		}()
	}
	uc.bus.PublishTopicEvent(topicSubscribeBidAskTargets, ctx, targetArr)
}

func tickProcessor(tickChan chan *entity.RealTimeTick) {
	for {
		tick := <-tickChan
		log.Infof("tick:%s\n", time.Since(tick.TickTime).String())
	}
}

func bidAskProcessor(bidAskChan chan *entity.RealTimeBidAsk) {
	for {
		bidAsk := <-bidAskChan
		log.Infof("bidask:%s\n", time.Since(bidAsk.TickTime).String())
	}
}
