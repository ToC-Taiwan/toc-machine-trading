package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/target"
	"tmt/internal/usecase/topic"
	"tmt/pkg/common"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo          TargetRepo
	gRPCAPI       TargetgRPCAPI
	streamgRPCAPI StreamgRPCAPI

	targetFilter      *target.Filter
	monitorFutureCode string
	waitMonitorFuture chan struct{}
}

// NewTarget -.
func NewTarget(r TargetRepo, t TargetgRPCAPI, s StreamgRPCAPI) Target {
	cfg := config.GetConfig()
	uc := &TargetUseCase{
		repo:              r,
		gRPCAPI:           t,
		streamgRPCAPI:     s,
		targetFilter:      target.NewFilter(cfg.TargetCond),
		waitMonitorFuture: make(chan struct{}),
	}

	bus.SubscribeTopic(topic.TopicMonitorFutureCode, uc.fillMonitorFutureCode)
	bus.PublishTopicEvent(topic.TopicQueryMonitorFutureCode)
	<-uc.waitMonitorFuture

	// unsubscriba all first
	if err := uc.UnSubscribeAll(context.Background()); err != nil {
		logger.Panic("unsubscribe all fail")
	}

	// query targets from db
	tradeDay := cc.GetBasicInfo().TradeDay
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), tradeDay)
	if err != nil {
		logger.Panic(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.SearchTradeDayTargets(context.Background(), tradeDay)
		if err != nil {
			logger.Panic(err)
		}

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			logger.Error("no targets")
			<-stuck
		}
	}

	uc.publishNewTargets(targetArr)

	// sub events
	bus.SubscribeTopic(topic.TopicNewTargets, uc.publishNewTargets)
	bus.SubscribeTopic(topic.TopicSubscribeStockTickTargets, uc.SubscribeStockTick, uc.SubscribeStockBidAsk)
	bus.SubscribeTopic(topic.TopicUnSubscribeStockTickTargets, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)
	bus.SubscribeTopic(topic.TopicSubscribeFutureTickTargets, uc.SubscribeFutureTick)

	return uc
}

func (uc *TargetUseCase) fillMonitorFutureCode(future *entity.Future) {
	uc.monitorFutureCode = future.Code
	close(uc.waitMonitorFuture)
}

func (uc *TargetUseCase) publishNewTargets(targetArr []*entity.StockTarget) {
	err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr)
	if err != nil {
		logger.Panic(err)
	}

	cc.AppendTargets(targetArr)

	// stock
	bus.PublishTopicEvent(topic.TopicFetchStockHistory, context.Background(), targetArr)
	bus.PublishTopicEvent(topic.TopicStreamStockTargets, context.Background(), targetArr)

	// future
	bus.PublishTopicEvent(topic.TopicStreamFutureTargets, context.Background(), uc.monitorFutureCode)
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.StockTarget {
	return cc.GetTargets()
}

// SearchTradeDayTargets - search targets by trade day
func (uc *TargetUseCase) SearchTradeDayTargets(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(common.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(cc.GetBasicInfo().TradeDay.Add(8*time.Hour)) {
		logger.Warn("VolumeRank is empty, search from all snapshot")
		return uc.SearchTradeDayTargetsFromAllSnapshot(tradeDay)
	}

	var result []*entity.StockTarget
	for _, v := range t {
		stock := cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.CheckVolume(v.GetTotalVolume()) || !uc.targetFilter.IsTarget(stock, v.GetClose()) {
			continue
		}

		result = append(result, &entity.StockTarget{
			Rank:     len(result) + 1,
			StockNum: v.GetCode(),
			Volume:   v.GetTotalVolume(),
			TradeDay: tradeDay,
			Stock:    stock,
		})
	}
	return result, nil
}

// SearchTradeDayTargetsFromAllSnapshot -.
func (uc *TargetUseCase) SearchTradeDayTargetsFromAllSnapshot(tradeDay time.Time) ([]*entity.StockTarget, error) {
	data, err := uc.streamgRPCAPI.GetAllStockSnapshot()
	if err != nil {
		return []*entity.StockTarget{}, err
	}

	if len(data) < 200 {
		return []*entity.StockTarget{}, errors.New("no all snapshots")
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})

	var result []*entity.StockTarget
	for _, v := range data[:200] {
		stock := cc.GetStockDetail(v.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.CheckVolume(v.GetTotalVolume()) || !uc.targetFilter.IsTarget(stock, v.GetClose()) {
			continue
		}

		result = append(result, &entity.StockTarget{
			Rank:     len(result) + 1,
			StockNum: v.GetCode(),
			Volume:   v.GetTotalVolume(),
			TradeDay: tradeDay,
			Stock:    stock,
		})
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
func (uc *TargetUseCase) SubscribeStockTick(targetArr []*entity.StockTarget) error {
	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
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
func (uc *TargetUseCase) SubscribeStockBidAsk(targetArr []*entity.StockTarget) error {
	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
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

// SubscribeFutureTick -.
func (uc *TargetUseCase) SubscribeFutureTick(code string) error {
	failSubNumArr, err := uc.gRPCAPI.SubscribeFutureTick([]string{code})
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe future fail %v", failSubNumArr)
	}

	return nil
}

// SubscribeFutureBidAsk -.
func (uc *TargetUseCase) SubscribeFutureBidAsk(code string) error {
	failSubNumArr, err := uc.gRPCAPI.SubscribeFutureBidAsk([]string{code})
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe future fail %v", failSubNumArr)
	}

	return nil
}
