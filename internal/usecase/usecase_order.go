package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/logger"
)

// OrderUseCase -.
type OrderUseCase struct {
	repo    OrderRepo
	gRPCAPI OrdergRPCAPI
	bus     *eventbus.Bus
}

// NewOrder -.
func NewOrder(r *repo.OrderRepo, t *grpcapi.OrdergRPCAPI, bus *eventbus.Bus) {
	uc := &OrderUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}

	if err := uc.bus.SubscribeTopic(topicTargets, uc.targetCallback); err != nil {
		logger.Get().Panic(err)
	}
}

func (uc *OrderUseCase) targetCallback(targetArr []*entity.Target) {
	for _, target := range targetArr {
		tickChan := make(chan *entity.RealTimeTick)
		bidAskChan := make(chan *entity.RealTimeBidAsk)

		CacheSetTickChan(target.StockNum, tickChan)
		CacheSetBidAskChan(target.StockNum, bidAskChan)

		go uc.tickProcessor(tickChan)
		go uc.bidAskProcessor(bidAskChan)
	}

	uc.bus.PublishTopicEvent(topicSubscribeTargets, context.Background(), targetArr)
}

func (uc *OrderUseCase) tickProcessor(tickChan chan *entity.RealTimeTick) {
	for {
		tick := <-tickChan
		logger.Get().Infof("tick:%s", tick.StockNum)
	}
}

func (uc *OrderUseCase) bidAskProcessor(bidAskChan chan *entity.RealTimeBidAsk) {
	for {
		bidAsk := <-bidAskChan
		logger.Get().Infof("bidask:%s", bidAsk.StockNum)
	}
}
