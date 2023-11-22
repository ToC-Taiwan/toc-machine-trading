package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"
	"tmt/internal/modules/target"
	"tmt/internal/modules/tradeday"
	"tmt/internal/usecase/grpc"
	"tmt/internal/usecase/repo"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

// TargetUseCase -.
type TargetUseCase struct {
	repo    TargetRepo
	gRPCAPI RealTimegRPCAPI

	targetFilter *target.Filter
	cfg          *config.Config
	tradeDay     *tradeday.TradeDay

	logger *log.Log
	cc     *Cache
	bus    *eventbus.Bus
}

func NewTarget() Target {
	cfg := config.Get()
	return &TargetUseCase{
		repo:         repo.NewTarget(cfg.GetPostgresPool()),
		gRPCAPI:      grpc.NewRealTime(cfg.GetSinopacPool()),
		cfg:          cfg,
		tradeDay:     tradeday.Get(),
		targetFilter: target.NewFilter(cfg.TargetStock),
	}
}

func (uc *TargetUseCase) Init(logger *log.Log, cc *Cache, bus *eventbus.Bus) Target {
	uc.logger = logger
	uc.cc = cc
	uc.bus = bus

	// query targets from db
	tDay := uc.tradeDay.GetStockTradeDay().TradeDay
	targetArr, err := uc.searchTradeDayTargetsFromDB(tDay)
	if err != nil {
		uc.logger.Fatal(err)
	}

	// db has no targets, find targets from gRPC
	if len(targetArr) == 0 {
		targetArr, err = uc.searchTradeDayTargets(tDay)
		if err != nil {
			uc.logger.Fatal(err)
		}

		if len(targetArr) == 0 {
			stuck := make(chan struct{})
			uc.logger.Error("no targets")
			<-stuck
		}
	}

	uc.cc.AppendStockTargets(targetArr)
	uc.publishNewStockTargets(targetArr)
	uc.publishNewFutureTargets()

	// go func() {
	// 	time.Sleep(time.Until(basic.TradeDay.Add(time.Hour * 9)))
	// 	for range time.NewTicker(time.Second * 60).C {
	// 		if uc.stockTradeInSwitch {
	// 			if err := uc.realTimeAddTargets(); err != nil {
	// 				uc.logger.Fatal(err)
	// 			}
	// 		}
	// 	}
	// }()

	return uc
}

// GetTargets - get targets from cache
func (uc *TargetUseCase) GetTargets(ctx context.Context) []*entity.StockTarget {
	return uc.cc.GetStockTargets()
}

func (uc *TargetUseCase) publishNewStockTargets(targetArr []*entity.StockTarget) {
	if err := uc.repo.InsertOrUpdateTargetArr(context.Background(), targetArr); err != nil {
		uc.logger.Fatal(err)
	}
	uc.bus.PublishTopicEvent(TopicFetchStockHistory, targetArr)
}

func (uc *TargetUseCase) publishNewFutureTargets() {
	if futureTarget, err := uc.getFutureTarget(); err != nil {
		uc.logger.Fatal(err)
	} else {
		uc.bus.PublishTopicEvent(TopicFetchFutureHistory, futureTarget)
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

func (uc *TargetUseCase) searchTradeDayTargetsFromDB(tradeDay time.Time) ([]*entity.StockTarget, error) {
	targetArr, err := uc.repo.QueryTargetsByTradeDay(context.Background(), tradeDay)
	if err != nil {
		return nil, err
	}

	var result []*entity.StockTarget
	for _, v := range targetArr {
		stock := uc.cc.GetStockDetail(v.StockNum)
		if stock == nil {
			continue
		}

		v.Stock = stock
		result = append(result, v)
	}

	return result, nil
}

func (uc *TargetUseCase) searchTradeDayTargets(tradeDay time.Time) ([]*entity.StockTarget, error) {
	lastTradeDay := uc.tradeDay.GetLastNStockTradeDay(1)[0]
	t, err := uc.gRPCAPI.GetStockVolumeRank(lastTradeDay.Format(global.ShortTimeLayout))
	if err != nil {
		return nil, err
	}

	if len(t) == 0 && time.Now().Before(uc.tradeDay.GetStockTradeDay().TradeDay.Add(8*time.Hour)) {
		uc.logger.Warn("VolumeRank is empty, search from all snapshot")
		return uc.searchTradeDayTargetsFromAllSnapshot(tradeDay)
	}

	var result []*entity.StockTarget
	for _, v := range t {
		stock := uc.cc.GetStockDetail(v.GetCode())
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
		stock := uc.cc.GetStockDetail(v.GetCode())
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
// 		uc.logger.Warnf("stock snapshot len is not enough: %d", len(data))
// 		return nil
// 	}

// 	sort.Slice(data, func(i, j int) bool {
// 		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
// 	})
// 	data = data[:uc.targetFilter.RealTimeRank]

// 	currentTargets := uc.cc.GetStockTargets()
// 	targetsMap := make(map[string]*entity.StockTarget)
// 	for _, t := range currentTargets {
// 		targetsMap[t.StockNum] = t
// 	}

// 	var newTargets []*entity.StockTarget
// 	for i, d := range data {
// 		stock := uc.cc.GetStockDetail(d.GetCode())
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
// 		uc.cc.AppendStockTargets(newTargets)
// 		uc.publishNewStockTargets(newTargets)
// 		for _, t := range newTargets {
// 			uc.logger.Infof("New target: %s", t.StockNum)
// 		}
// 	}
// 	return nil
// }
