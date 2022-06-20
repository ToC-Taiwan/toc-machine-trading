package usecase

import (
	"context"
	"errors"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/global"
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
		log.Panic("unsubscribe all fail")
	}

	tradeDay := cc.GetBasicInfo().TradeDay
	targetArr, err := uc.repo.QueryTargetsByTradeDay(ctx, tradeDay)
	if err != nil {
		log.Panic(err)
	}

	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(ctx, tradeDay)
		if err != nil {
			log.Panic(err)
		}

		if len(targetArr) != 0 {
			if err = uc.repo.InsertTargetArr(ctx, targetArr); err != nil {
				log.Panic(err)
			}
		}
	}

	// pub events
	bus.PublishTopicEvent(topicTargets, ctx, targetArr)

	// sub events
	if err := bus.SubscribeTopic(topicSubscribeTickTargets, uc.SubscribeStockTick, uc.SubscribeStockBidAsk); err != nil {
		log.Panic(err)
	}
}

// SearchTradeDayTargets -.
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	cond := cfg.TargetCond
	var result []*entity.Target
	for i, v := range t {
		if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh || v.GetTotalAmount() < cond.LimitVolume {
			continue
		}
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
