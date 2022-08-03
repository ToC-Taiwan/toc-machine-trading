package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/global"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI TargetgRPCAPI
}

// NewTarget -.
func NewTarget(r TargetRepo, t TargetgRPCAPI) *TargetUseCase {
	uc := &TargetUseCase{
		repo:    r,
		gRPCAPI: t,
	}

	// unsubscriba all first
	if err := uc.UnSubscribeAll(context.Background()); err != nil {
		log.Panic("unsubscribe all fail")
	}

	// query targets from db
	tradeDay := cc.GetBasicInfo().TradeDay
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), tradeDay)
	if err != nil {
		log.Panic(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(context.Background(), tradeDay)
		if err != nil {
			log.Panic(err)
		}

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			log.Error("no targets")
			<-stuck
		}
	}

	uc.publishNewTargets(targetArr, false)

	// sub events
	bus.SubscribeTopic(topicRealTimeTargets, uc.publishNewTargets)
	bus.SubscribeTopic(topicSubscribeTickTargets, uc.SubscribeStockTick, uc.SubscribeStockBidAsk)
	bus.SubscribeTopic(topicUnSubscribeTickTargets, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)

	return uc
}

func (uc *TargetUseCase) publishNewTargets(targetArr []*entity.Target, subscribe bool) {
	err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr)
	if err != nil {
		log.Panic(err)
	}

	bus.PublishTopicEvent(topicFetchHistory, context.Background(), targetArr)
	if subscribe {
		bus.PublishTopicEvent(topicStreamTargets, context.Background(), targetArr)
	}
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.Target {
	return cc.GetTargets()
}

// InsertTargets - insert targets to db
func (uc *TargetUseCase) InsertTargets(ctx context.Context, targetArr []*entity.Target) error {
	if err := uc.repo.InsertTargetArr(ctx, targetArr); err != nil {
		return err
	}
	return nil
}

// SearchTradeDayTargets - search targets by trade day
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
	condition := cfg.TargetCond

	var result []*entity.Target
	for _, c := range condition.PriceVolumeLimit {
		for _, v := range t {
			stock := cc.GetStockDetail(v.GetCode())
			if stock == nil {
				continue
			}

			if !blackStockFilter(stock.Number, condition) || !blackCatagoryFilter(stock.Category, condition) {
				continue
			}

			if !targetFilter(v.GetClose(), v.GetTotalVolume(), c, false) {
				continue
			}

			result = append(result, &entity.Target{
				Rank:     len(result) + 1,
				StockNum: v.GetCode(),
				Volume:   v.GetTotalVolume(),
				PreFetch: c.PreFetch,
				RealTime: false,
				TradeDay: tradeDay,
				Stock:    stock,
			})
		}
	}
	return result, nil
}

// UnSubscribeAll -.
func (uc *TargetUseCase) UnSubscribeAll(ctx context.Context) error {
	failMessge, err := uc.gRPCAPI.UnSubscribeStockAllTick()
	if err != nil {
		return err
	}

	if m := failMessge.GetErr(); m != "" {
		return errors.New(m)
	}

	failMessge, err = uc.gRPCAPI.UnSubscribeStockAllBidAsk()
	if err != nil {
		return err
	}

	if m := failMessge.GetErr(); m != "" {
		return errors.New(m)
	}

	return nil
}

// SubscribeStockTick -.
func (uc *TargetUseCase) SubscribeStockTick(targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		if v.RealTime {
			subArr = append(subArr, v.StockNum)
		}
	}

	failSubNumArr, err := uc.gRPCAPI.SubscribeStockTick(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// SubscribeStockBidAsk -.
func (uc *TargetUseCase) SubscribeStockBidAsk(targetArr []*entity.Target) error {
	var subArr []string
	for _, v := range targetArr {
		if v.RealTime {
			subArr = append(subArr, v.StockNum)
		}
	}

	failSubNumArr, err := uc.gRPCAPI.SubscribeStockBidAsk(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// UnSubscribeStockTick -.
func (uc *TargetUseCase) UnSubscribeStockTick(stockNum string) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockTick([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}

// UnSubscribeStockBidAsk -.
func (uc *TargetUseCase) UnSubscribeStockBidAsk(stockNum string) error {
	failUnSubNumArr, err := uc.gRPCAPI.UnSubscribeStockBidAsk([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}
