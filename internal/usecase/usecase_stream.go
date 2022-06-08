package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/logger"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	gRPCAPI StreamgRPCAPI
}

// NewStream -.
func NewStream(r *repo.StreamRepo, t *grpcapi.StreamgRPCAPI) {
	uc := &StreamUseCase{
		repo:    r,
		gRPCAPI: t,
	}

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveTicks(context.Background())
	go uc.ReceiveBidAsk(context.Background())
	go uc.ReceiveOrderStatus(context.Background())
}

// ReceiveEvent -.
func (uc *StreamUseCase) ReceiveEvent(ctx context.Context) {
	eventChan := make(chan *entity.SinopacEvent)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				logger.Get().Warn(event)
				if err := uc.repo.InsertEvent(ctx, event); err != nil {
					logger.Get().Panic(err)
				}
			}
		}
	}()

	if err := uc.gRPCAPI.EventChannel(eventChan); err != nil {
		logger.Get().Panic(err)
	}
}

// ReceiveTicks -.
func (uc *StreamUseCase) ReceiveTicks(ctx context.Context) {
	tickChan := make(chan *entity.RealTimeTick)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tick := <-tickChan:
				CacheGetTickChan(tick.StockNum) <- tick
			}
		}
	}()

	if err := uc.gRPCAPI.TickChannel(tickChan); err != nil {
		logger.Get().Panic(err)
	}
}

// ReceiveBidAsk -.
func (uc *StreamUseCase) ReceiveBidAsk(ctx context.Context) {
	bidAskChan := make(chan *entity.RealTimeBidAsk)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case bidAsk := <-bidAskChan:
				CacheGetBidAskChan(bidAsk.StockNum) <- bidAsk
			}
		}
	}()

	if err := uc.gRPCAPI.BidAskChannel(bidAskChan); err != nil {
		logger.Get().Panic(err)
	}
}

// ReceiveOrderStatus -.
func (uc *StreamUseCase) ReceiveOrderStatus(ctx context.Context) {
	orderStatusChan := make(chan *entity.OrderStatus)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case orderStatus := <-orderStatusChan:
				logger.Get().Info(orderStatus)
			}
		}
	}()

	if err := uc.gRPCAPI.OrderStatusChannel(orderStatusChan); err != nil {
		logger.Get().Panic(err)
	}
}
