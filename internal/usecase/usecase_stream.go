package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"
	"tmt/internal/usecase/modules/target"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/internal/usecase/modules/trader"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	rabbit  StreamRabbit
	grpcapi StreamgRPCAPI

	basic entity.BasicInfo

	stockTradeSwitchCfg  config.StockTradeSwitch
	futureTradeSwitchCfg config.FutureTradeSwitch

	stockAnalyzeCfg  config.StockAnalyze
	futureAnalyzeCfg config.FutureAnalyze

	targetFilter *target.Filter
	tradeDay     *tradeday.TradeDay

	stockTradeInSwitch  bool
	futureTradeInSwitch bool

	mainFutureCode string
	tradeIndex     *entity.TradeIndex
}

// NewStream -.
func NewStream(r StreamRepo, g StreamgRPCAPI, t StreamRabbit) *StreamUseCase {
	cfg := config.GetConfig()
	uc := &StreamUseCase{
		repo:    r,
		rabbit:  t,
		grpcapi: g,

		stockTradeSwitchCfg:  cfg.StockTradeSwitch,
		futureTradeSwitchCfg: cfg.FutureTradeSwitch,
		stockAnalyzeCfg:      cfg.StockAnalyze,
		futureAnalyzeCfg:     cfg.FutureAnalyze,

		basic:        *cc.GetBasicInfo(),
		targetFilter: target.NewFilter(cfg.TargetCond),
		tradeDay:     tradeday.NewTradeDay(),
	}
	t.FillAllBasic(uc.basic.AllStocks, uc.basic.AllFutures)

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())
	go uc.periodUpdateTradeIndex()

	go func() {
		time.Sleep(time.Until(cc.GetBasicInfo().TradeDay.Add(time.Hour * 9)))
		for range time.NewTicker(time.Second * 60).C {
			if uc.stockTradeInSwitch {
				if err := uc.realTimeAddTargets(); err != nil {
					log.Panic(err)
				}
			}
		}
	}()

	bus.SubscribeTopic(event.TopicStreamStockTargets, uc.ReceiveStreamData)
	bus.SubscribeTopic(event.TopicStreamFutureTargets, uc.ReceiveFutureStreamData)
	bus.SubscribeTopic(event.TopicMonitorFutureCode, uc.updateMainFutureCode)
	return uc
}

func (uc *StreamUseCase) GetTradeIndex() *entity.TradeIndex {
	return uc.tradeIndex
}

func (uc *StreamUseCase) periodUpdateTradeIndex() {
	for range time.NewTicker(time.Second * 5).C {
		if uc.tradeIndex == nil {
			uc.tradeIndex = &entity.TradeIndex{}
		}

		var err error
		ctx := context.Background()
		go func() {
			uc.tradeIndex.Nasdaq, err = uc.GetNasdaqClose()
			if err != nil {
				log.Error(err)
			}
		}()

		go func() {
			uc.tradeIndex.TSE, err = uc.GetTSESnapshot(ctx)
			if err != nil {
				log.Error(err)
			}
		}()

		go func() {
			uc.tradeIndex.NF, err = uc.GetNasdaqFutureClose()
			if err != nil {
				log.Error(err)
			}
		}()

		go func() {
			uc.tradeIndex.OTC, err = uc.GetOTCSnapshot(ctx)
			if err != nil {
				log.Error(err)
			}
		}()
	}
}

func (uc *StreamUseCase) realTimeAddTargets() error {
	data, err := uc.grpcapi.GetAllStockSnapshot()
	if err != nil {
		return err
	}

	// at least 200 snapshot to rank volume
	if len(data) < 200 {
		log.Warnf("stock snapshot len is not enough: %d", len(data))
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:uc.targetFilter.RealTimeRank]

	currentTargets := cc.GetTargets()
	targetsMap := make(map[string]*entity.StockTarget)
	for _, t := range currentTargets {
		targetsMap[t.StockNum] = t
	}

	var newTargets []*entity.StockTarget
	for i, d := range data {
		stock := cc.GetStockDetail(d.GetCode())
		if stock == nil {
			continue
		}

		if !uc.targetFilter.IsTarget(stock, d.GetClose()) {
			continue
		}

		if targetsMap[d.GetCode()] == nil {
			newTargets = append(newTargets, &entity.StockTarget{
				Rank:     100 + i + 1,
				StockNum: d.GetCode(),
				Volume:   d.GetTotalVolume(),
				TradeDay: uc.basic.TradeDay,
				Stock:    stock,
			})
		}
	}

	if len(newTargets) != 0 {
		bus.PublishTopicEvent(event.TopicNewTargets, newTargets)
	}
	return nil
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

			if event.EventCode != 16 {
				log.Warnf("EventCode: %d, Event: %s, ResoCode: %d, Info: %s", event.EventCode, event.Event, event.Response, event.Info)
			}
		}
	}()
	uc.rabbit.EventConsumer(eventChan)
}

// ReceiveOrderStatus -.
func (uc *StreamUseCase) ReceiveOrderStatus(ctx context.Context) {
	orderStatusChan := make(chan interface{})
	go func() {
		for {
			order := <-orderStatusChan
			switch t := order.(type) {
			case *entity.StockOrder:
				if cc.GetOrderByOrderID(t.OrderID) == nil {
					cc.SetOrderByOrderID(t.ToManual())
				}
				bus.PublishTopicEvent(event.TopicInsertOrUpdateStockOrder, t)
			case *entity.FutureOrder:
				if cc.GetFutureOrderByOrderID(t.OrderID) == nil {
					cc.SetFutureOrderByOrderID(t.ToManual())
				}
				bus.PublishTopicEvent(event.TopicInsertOrUpdateFutureOrder, t)
			}
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// NewOrderStatusConnection -.
func (uc *StreamUseCase) NewOrderStatusConnection(orderStatusChan chan interface{}, connectionID string) {
	uc.rabbit.AddOrderStatusChan(orderStatusChan, connectionID)
}

func (uc *StreamUseCase) DeleteOrderStatusConnection(connectionID string) {
	uc.rabbit.RemoveOrderStatusChan(connectionID)
}

// GetTSESnapshot -.
func (uc *StreamUseCase) GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
	body, err := uc.grpcapi.GetStockSnapshotTSE()
	if err != nil {
		return nil, err
	}
	return &entity.StockSnapShot{
		SnapShotBase: entity.SnapShotBase{
			SnapTime:        time.Unix(0, body.GetTs()).Add(-8 * time.Hour),
			Open:            body.GetOpen(),
			High:            body.GetHigh(),
			Low:             body.GetLow(),
			Close:           body.GetClose(),
			TickType:        body.GetTickType(),
			PriceChg:        body.GetChangePrice(),
			PctChg:          body.GetChangeRate(),
			ChgType:         body.GetChangeType(),
			Volume:          body.GetVolume(),
			VolumeSum:       body.GetTotalVolume(),
			Amount:          body.GetAmount(),
			AmountSum:       body.GetTotalAmount(),
			YesterdayVolume: body.GetYesterdayVolume(),
			VolumeRatio:     body.GetVolumeRatio(),
		},
		StockNum: body.GetCode(),
	}, nil
}

// GetOTCSnapshot -.
func (uc *StreamUseCase) GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
	body, err := uc.grpcapi.GetStockSnapshotOTC()
	if err != nil {
		return nil, err
	}
	return &entity.StockSnapShot{
		SnapShotBase: entity.SnapShotBase{
			SnapTime:        time.Unix(0, body.GetTs()).Add(-8 * time.Hour),
			Open:            body.GetOpen(),
			High:            body.GetHigh(),
			Low:             body.GetLow(),
			Close:           body.GetClose(),
			TickType:        body.GetTickType(),
			PriceChg:        body.GetChangePrice(),
			PctChg:          body.GetChangeRate(),
			ChgType:         body.GetChangeType(),
			Volume:          body.GetVolume(),
			VolumeSum:       body.GetTotalVolume(),
			Amount:          body.GetAmount(),
			AmountSum:       body.GetTotalAmount(),
			YesterdayVolume: body.GetYesterdayVolume(),
			VolumeRatio:     body.GetVolumeRatio(),
		},
		StockNum: body.GetCode(),
	}, nil
}

func (uc *StreamUseCase) GetNasdaqClose() (*entity.YahooPrice, error) {
	d, err := uc.grpcapi.GetNasdaq()
	if err != nil {
		return nil, err
	}

	return &entity.YahooPrice{
		Last:      d.GetLast(),
		Price:     d.GetPrice(),
		UpdatedAt: time.Now(),
	}, nil
}

func (uc *StreamUseCase) GetNasdaqFutureClose() (*entity.YahooPrice, error) {
	d, err := uc.grpcapi.GetNasdaqFuture()
	if err != nil {
		return nil, err
	}
	return &entity.YahooPrice{
		Last:      d.GetLast(),
		Price:     d.GetPrice(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetStockSnapshotByNumArr -.
func (uc *StreamUseCase) GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error) {
	var fetchArr, stockNotExist []string
	for _, s := range stockNumArr {
		if cc.GetStockDetail(s) == nil {
			stockNotExist = append(stockNotExist, s)
		} else {
			fetchArr = append(fetchArr, s)
		}
	}

	snapshot, err := uc.grpcapi.GetStockSnapshotByNumArr(fetchArr)
	if err != nil {
		return nil, err
	}

	var result []*entity.StockSnapShot
	for _, body := range snapshot {
		stockNum := body.GetCode()
		result = append(result, &entity.StockSnapShot{
			StockNum:  stockNum,
			StockName: cc.GetStockDetail(stockNum).Name,
			SnapShotBase: entity.SnapShotBase{
				SnapTime:        time.Unix(0, body.GetTs()).Add(-8 * time.Hour),
				Open:            body.GetOpen(),
				High:            body.GetHigh(),
				Low:             body.GetLow(),
				Close:           body.GetClose(),
				TickType:        body.GetTickType(),
				PriceChg:        body.GetChangePrice(),
				PctChg:          body.GetChangeRate(),
				ChgType:         body.GetChangeType(),
				Volume:          body.GetVolume(),
				VolumeSum:       body.GetTotalVolume(),
				Amount:          body.GetAmount(),
				AmountSum:       body.GetTotalAmount(),
				YesterdayVolume: body.GetYesterdayVolume(),
				VolumeRatio:     body.GetVolumeRatio(),
			},
		})
	}

	for _, v := range stockNotExist {
		result = append(result, &entity.StockSnapShot{
			StockNum: v,
		})
	}

	return result, nil
}

// GetFutureSnapshotByCode -.
func (uc *StreamUseCase) GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error) {
	snapshot, err := uc.grpcapi.GetFutureSnapshotByCode(code)
	if err != nil {
		return nil, err
	}

	return &entity.FutureSnapShot{
		Code:       snapshot.GetCode(),
		FutureName: cc.GetFutureDetail(code).Name,
		SnapShotBase: entity.SnapShotBase{
			SnapTime:        time.Unix(0, snapshot.GetTs()).Add(-8 * time.Hour),
			Open:            snapshot.GetOpen(),
			High:            snapshot.GetHigh(),
			Low:             snapshot.GetLow(),
			Close:           snapshot.GetClose(),
			TickType:        snapshot.GetTickType(),
			PriceChg:        snapshot.GetChangePrice(),
			PctChg:          snapshot.GetChangeRate(),
			ChgType:         snapshot.GetChangeType(),
			Volume:          snapshot.GetVolume(),
			VolumeSum:       snapshot.GetTotalVolume(),
			Amount:          snapshot.GetAmount(),
			AmountSum:       snapshot.GetTotalAmount(),
			YesterdayVolume: snapshot.GetYesterdayVolume(),
			VolumeRatio:     snapshot.GetVolumeRatio(),
		},
	}, nil
}

// NewFutureRealTimeConnection -.
func (uc *StreamUseCase) NewFutureRealTimeConnection(tickChan chan *entity.RealTimeFutureTick, connectionID string) {
	uc.rabbit.AddFutureTickChan(tickChan, connectionID)
}

// DeleteFutureRealTimeConnection -.
func (uc *StreamUseCase) DeleteFutureRealTimeConnection(connectionID string) {
	uc.rabbit.RemoveFutureTickChan(connectionID)
}

func (uc *StreamUseCase) updateMainFutureCode(future *entity.Future) {
	uc.mainFutureCode = future.Code
}

// GetMainFutureCode -.
func (uc *StreamUseCase) GetMainFutureCode() string {
	return uc.mainFutureCode
}

// ReceiveStreamData - receive target data, start goroutine to trade
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.StockTarget) {
	agentChan := make(chan *trader.StockTrader)
	targetMap := make(map[string]*entity.StockTarget)
	mutex := sync.RWMutex{}

	go func() {
		for {
			agent, ok := <-agentChan
			if !ok {
				break
			}
			go agent.TradingRoom()

			// send tick, bidask to trade room's channel
			go uc.rabbit.TickConsumer(agent.GetStockNum(), agent.GetTickChan())
			go uc.rabbit.StockBidAskConsumer(agent.GetStockNum(), agent.GetBidAskChan())

			mutex.RLock()
			target := targetMap[agent.GetStockNum()]
			mutex.RUnlock()

			bus.PublishTopicEvent(event.TopicSubscribeStockTickTargets, []*entity.StockTarget{target})
		}
	}()

	var wg sync.WaitGroup
	for _, t := range targetArr {
		target := t
		mutex.Lock()
		targetMap[target.StockNum] = target
		mutex.Unlock()

		wg.Add(1)
		go func() {
			defer wg.Done()
			agent := trader.NewStockTrader(target.StockNum, uc.stockTradeSwitchCfg, uc.stockAnalyzeCfg)
			agentChan <- agent
		}()
	}

	wg.Wait()
	close(agentChan)

	go uc.checkStockTradeSwitch()
}

// ReceiveFutureStreamData -.
func (uc *StreamUseCase) ReceiveFutureStreamData(ctx context.Context, code string) {
	agent := trader.NewFutureTrader(code, uc.futureTradeSwitchCfg, uc.futureAnalyzeCfg)

	go agent.TradingRoom()

	go uc.rabbit.FutureTickConsumer(code, agent.GetTickChan())
	// go uc.rabbit.FutureBidAskConsumer(code, agent.GetBidAskChan())

	bus.PublishTopicEvent(event.TopicSubscribeFutureTickTargets, code)

	go uc.checkFutureTradeSwitch()
}

func (uc *StreamUseCase) checkStockTradeSwitch() {
	if !uc.stockTradeSwitchCfg.AllowTrade {
		return
	}

	openTime := uc.basic.OpenTime
	tradeInEndTime := uc.basic.TradeInEndTime

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		var tempSwitch bool
		switch {
		case now.Before(openTime) || now.After(tradeInEndTime):
			tempSwitch = false
		case now.After(openTime) && now.Before(tradeInEndTime):
			tempSwitch = true
		}

		if uc.stockTradeInSwitch != tempSwitch {
			uc.stockTradeInSwitch = tempSwitch
			bus.PublishTopicEvent(event.TopicUpdateStockTradeSwitch, uc.stockTradeInSwitch)
		}
	}
}

func (uc *StreamUseCase) checkFutureTradeSwitch() {
	if !uc.futureTradeSwitchCfg.AllowTrade {
		return
	}

	futureTradeDay := uc.tradeDay.GetFutureTradeDay()
	timeRange := [][]time.Time{}
	firstStart := futureTradeDay.StartTime
	secondStart := futureTradeDay.EndTime.Add(-300 * time.Minute)

	timeRange = append(timeRange, []time.Time{firstStart, firstStart.Add(time.Duration(uc.futureTradeSwitchCfg.TradeTimeRange.FirstPartDuration) * time.Minute)})
	timeRange = append(timeRange, []time.Time{secondStart, secondStart.Add(time.Duration(uc.futureTradeSwitchCfg.TradeTimeRange.SecondPartDuration) * time.Minute)})

	for range time.NewTicker(2500 * time.Millisecond).C {
		now := time.Now()
		var tempSwitch bool
		for _, rangeTime := range timeRange {
			if now.After(rangeTime[0]) && now.Before(rangeTime[1]) {
				tempSwitch = true
			}
		}

		if uc.futureTradeInSwitch != tempSwitch {
			uc.futureTradeInSwitch = tempSwitch
			bus.PublishTopicEvent(event.TopicUpdateFutureTradeSwitch, uc.futureTradeInSwitch)
		}
	}
}

func (uc *StreamUseCase) TurnFutureTradeSwitch(ctx context.Context, allow bool) {
	uc.futureTradeSwitchCfg.AllowTrade = allow
}

func (uc *StreamUseCase) GetFutureTradeSwitchStatus(ctx context.Context) bool {
	return uc.futureTradeSwitchCfg.AllowTrade
}
