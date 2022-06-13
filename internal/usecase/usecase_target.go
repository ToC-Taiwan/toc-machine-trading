package usecase

import (
	"context"
	"errors"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/global"
	"toc-machine-trading/pkg/logger"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI TargetRPCAPI
	bus     *eventbus.Bus
}

// NewTarget -.
func NewTarget(r *repo.TargetRepo, t *grpcapi.TargetgRPCAPI, bus *eventbus.Bus) {
	uc := &TargetUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}

	ctx := context.Background()

	// unsubscriba all first
	if err := uc.UnSubscribeAll(ctx); err != nil {
		logger.Get().Panic("unsubscribe all fail")
	}

	targetArr, err := uc.repo.QueryTargetsByTradeDay(ctx, CacheGetTradeDay())
	if err != nil {
		logger.Get().Panic(err)
	}

	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(ctx, CacheGetTradeDay())
		if err != nil {
			logger.Get().Panic(err)
		}

		if len(targetArr) != 0 {
			if err = uc.repo.InsertTargetArr(ctx, targetArr); err != nil {
				logger.Get().Panic(err)
			}
		}
	}

	// pub events
	bus.PublishTopicEvent(topicTargets, ctx, targetArr)

	// sub events
	if err := bus.SubscribeTopic(topicSubscribeTickTargets, uc.SubscribeStockTick); err != nil {
		logger.Get().Panic(err)
	}
	if err := bus.SubscribeTopic(topicSubscribeBidAskTargets, uc.SubscribeStockBidAsk); err != nil {
		logger.Get().Panic(err)
	}
}

// SearchTradeDayTargets -.
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error) {
	lastTradeDay := GetLastNTradeDayByDate(1, tradeDay)[0]
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}
	var result []*entity.Target
	for i, v := range t {
		result = append(result, &entity.Target{
			StockNum:    v.GetCode(),
			TradeDay:    tradeDay,
			Rank:        i + 1,
			Volume:      v.GetTotalVolume(),
			Subscribe:   true,
			RealTimeAdd: false,
		})
	}
	return result, nil
}

// UnSubscribeAll -.
func (uc *TargetUseCase) UnSubscribeAll(ctx context.Context) error {
	fail, err := uc.gRPCAPI.UnSubscribeStockAllTick()
	if err != nil {
		return err
	}

	if fail.GetErr() != "" {
		return errors.New(fail.GetErr())
	}

	fail, err = uc.gRPCAPI.UnSubscribeStockAllBidAsk()
	if err != nil {
		return err
	}

	if fail.GetErr() != "" {
		return errors.New(fail.GetErr())
	}

	return nil
}

// SubscribeStockTick -.
func (uc *TargetUseCase) SubscribeStockTick(ctx context.Context, targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
	}

	fail, err := uc.gRPCAPI.SubscribeStockTick(subArr)
	if err != nil {
		return err
	}

	if len(fail) != 0 {
		return errors.New("subscribe fail")
	}

	return nil
}

// SubscribeStockBidAsk -.
func (uc *TargetUseCase) SubscribeStockBidAsk(ctx context.Context, targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
	}

	fail, err := uc.gRPCAPI.SubscribeStockBidAsk(subArr)
	if err != nil {
		return err
	}

	if len(fail) != 0 {
		return errors.New("subscribe fail")
	}

	return nil
}
