package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/event"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/target"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/repo"
	"tmt/pkg/common"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI RealTimegRPCAPI

	targetFilter *target.Filter
	cfg          *config.Config
	basic        *entity.BasicInfo
	tradeDay     *tradeday.TradeDay
}

func (u *UseCaseBase) NewTarget() Target {
	cfg := u.cfg
	basic := cc.GetBasicInfo()
	uc := &TargetUseCase{
		repo:         repo.NewTarget(u.pg),
		gRPCAPI:      grpcapi.NewRealTime(u.sc),
		cfg:          cfg,
		basic:        basic,
		tradeDay:     tradeday.Get(),
		targetFilter: target.NewFilter(cfg.TargetStock),
	}

	// query targets from db
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), uc.basic.TradeDay)
	if err != nil {
		logger.Fatal(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.searchTradeDayTargets(uc.basic.TradeDay)
		if err != nil {
			logger.Fatal(err)
		}

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			logger.Error("no targets")
			<-stuck
		}
	}

	cc.AppendStockTargets(targetArr)
	uc.publishNewStockTargets(targetArr)
	uc.publishNewFutureTargets()

	// go func() {
	// 	time.Sleep(time.Until(basic.TradeDay.Add(time.Hour * 9)))
	// 	for range time.NewTicker(time.Second * 60).C {
	// 		if uc.stockTradeInSwitch {
	// 			if err := uc.realTimeAddTargets(); err != nil {
	// 				logger.Fatal(err)
	// 			}
	// 		}
	// 	}
	// }()

	return uc
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.StockTarget {
	return cc.GetStockTargets()
}

func (uc *TargetUseCase) publishNewStockTargets(targetArr []*entity.StockTarget) {
	if err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr); err != nil {
		logger.Fatal(err)
	}
	bus.PublishTopicEvent(event.TopicFetchStockHistory, targetArr)
}

func (uc *TargetUseCase) publishNewFutureTargets() {
	if futureTarget, err := uc.getFutureTarget(); err != nil {
		logger.Fatal(err)
	} else {
		bus.PublishTopicEvent(event.TopicSubscribeFutureTickTargets, futureTarget)
	}
}

func (uc *TargetUseCase) getFutureTarget() (string, error) {
	futures, err := uc.repo.QueryAllMXFFuture(context.Background())
	if err != nil {
		return "", err
	}

	for _, v := range futures {
		if v.Code == "MXFR1" || v.Code == "MXFR2" {
			continue
		}

		if time.Now().Before(v.DeliveryDate) {
			return v.Code, nil
		}
	}

	return "", errors.New("no future")
}

func (uc *TargetUseCase) searchTradeDayTargets(tradeDay time.Time) ([]*entity.StockTarget, error) {
	lastTradeDay := cc.GetBasicInfo().LastTradeDay
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(common.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(cc.GetBasicInfo().TradeDay.Add(8*time.Hour)) {
		logger.Warn("VolumeRank is empty, search from all snapshot")
		return uc.searchTradeDayTargetsFromAllSnapshot(tradeDay)
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

func (uc *TargetUseCase) searchTradeDayTargetsFromAllSnapshot(tradeDay time.Time) ([]*entity.StockTarget, error) {
	data, err := uc.gRPCAPI.GetAllStockSnapshot()
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

// func (uc *TargetUseCase) realTimeAddTargets() error {
// 	data, err := uc.gRPCAPI.GetAllStockSnapshot()
// 	if err != nil {
// 		return err
// 	}

// 	// at least 200 snapshot to rank volume
// 	if len(data) < 200 {
// 		logger.Warnf("stock snapshot len is not enough: %d", len(data))
// 		return nil
// 	}

// 	sort.Slice(data, func(i, j int) bool {
// 		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
// 	})
// 	data = data[:uc.targetFilter.RealTimeRank]

// 	currentTargets := cc.GetStockTargets()
// 	targetsMap := make(map[string]*entity.StockTarget)
// 	for _, t := range currentTargets {
// 		targetsMap[t.StockNum] = t
// 	}

// 	var newTargets []*entity.StockTarget
// 	for i, d := range data {
// 		stock := cc.GetStockDetail(d.GetCode())
// 		if stock == nil {
// 			continue
// 		}

// 		if !uc.targetFilter.IsTarget(stock, d.GetClose()) {
// 			continue
// 		}

// 		if targetsMap[d.GetCode()] == nil {
// 			newTargets = append(newTargets, &entity.StockTarget{
// 				Rank:     100 + i + 1,
// 				StockNum: d.GetCode(),
// 				Volume:   d.GetTotalVolume(),
// 				TradeDay: uc.basic.TradeDay,
// 				Stock:    stock,
// 			})
// 		}
// 	}

// 	if len(newTargets) != 0 {
// 		cc.AppendStockTargets(newTargets)
// 		uc.publishNewStockTargets(newTargets)
// 		for _, t := range newTargets {
// 			logger.Infof("New target: %s", t.StockNum)
// 		}
// 	}
// 	return nil
// }
