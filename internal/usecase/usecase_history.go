package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
)

// HistoryUseCase -.
type HistoryUseCase struct {
	repo    HistoryRepo
	grpcapi HistorygRPCAPI
	bus     *eventbus.Bus
}

// NewHistory -.
func NewHistory(r *repo.HistoryRepo, t *grpcapi.HistorygRPCAPI, bus *eventbus.Bus) {
	uc := &HistoryUseCase{
		repo:    r,
		grpcapi: t,
		bus:     bus,
	}

	if err := uc.bus.SubscribeTopic(topicTargets, uc.FetchHistory); err != nil {
		log.Panic(err)
	}
}

// FetchHistory FetchHistory
func (uc *HistoryUseCase) FetchHistory(ctx context.Context, targetArr []*entity.Target) {
	uc.bus.PublishTopicEvent(topicStreamTickTargets, ctx, targetArr)
	uc.bus.PublishTopicEvent(topicStreamBidAskTargets, ctx, targetArr)
}
