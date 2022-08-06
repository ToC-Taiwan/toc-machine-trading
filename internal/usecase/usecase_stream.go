package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
)

// StreamUseCase -.
type StreamUseCase struct {
	repo    StreamRepo
	rabbit  StreamRabbit
	grpcapi StreamgRPCAPI

	tradeSwitchCfg config.TradeSwitch
	analyzeCfg     config.Analyze
	basic          entity.BasicInfo
	targetCond     config.TargetCond

	tradeInSwitch bool
	clearAll      bool
}

// NewStream -.
func NewStream(r StreamRepo, g StreamgRPCAPI, t StreamRabbit) *StreamUseCase {
	uc := &StreamUseCase{
		repo:    r,
		rabbit:  t,
		grpcapi: g,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	uc.tradeSwitchCfg = cfg.TradeSwitch
	uc.analyzeCfg = cfg.Analyze
	uc.targetCond = cfg.TargetCond
	uc.basic = *cc.GetBasicInfo()

	go uc.checkTradeSwitch()
	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	go func() {
		time.Sleep(time.Until(cc.GetBasicInfo().TradeDay.Add(time.Hour * 9)))
		for range time.NewTicker(time.Second * 20).C {
			if uc.tradeInSwitch {
				if err := uc.realTimeAddTargets(); err != nil {
					log.Panic(err)
				}
			}
		}
	}()

	bus.SubscribeTopic(topicStreamTargets, uc.ReceiveStreamData)
	return uc
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
	orderStatusChan := make(chan *entity.Order)
	go func() {
		for {
			order := <-orderStatusChan
			bus.PublishTopicEvent(topicInsertOrUpdateOrder, order)
		}
	}()
	uc.rabbit.OrderStatusConsumer(orderStatusChan)
}

// ReceiveStreamData - receive target data, start goroutine to trade
func (uc *StreamUseCase) ReceiveStreamData(ctx context.Context, targetArr []*entity.Target) {
	agentChan := make(chan *TradeAgent)
	targetMap := make(map[string]*entity.Target)
	mutex := sync.RWMutex{}

	go func() {
		for {
			agent, ok := <-agentChan
			if !ok {
				break
			}
			go uc.tradingRoom(agent)

			// send tick, bidask to trade room's channel
			go uc.rabbit.TickConsumer(agent.stockNum, agent.tickChan)
			go uc.rabbit.BidAskConsumer(agent.stockNum, agent.bidAskChan)

			mutex.RLock()
			target := targetMap[agent.stockNum]
			mutex.RUnlock()

			bus.PublishTopicEvent(topicSubscribeTickTargets, []*entity.Target{target})
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
			agent := NewAgent(target.StockNum, uc.tradeSwitchCfg)
			agentChan <- agent
		}()
	}
	wg.Wait()
	close(agentChan)
}

func (uc *StreamUseCase) tradingRoom(agent *TradeAgent) {
	go func() {
		for {
			agent.lastTick = <-agent.tickChan
			agent.tickArr = append(agent.tickArr, agent.lastTick)

			log.Debugf("%s tick time delay: %s", agent.stockNum, time.Since(agent.lastTick.TickTime).String())

			order := agent.generateOrder(uc.analyzeCfg, uc.clearAll)
			if order == nil {
				continue
			}

			uc.placeOrder(agent, order)
		}
	}()

	go func() {
		for {
			agent.lastBidAsk = <-agent.bidAskChan
			log.Debugf("%s bidask time delay: %s", agent.stockNum, time.Since(agent.lastBidAsk.BidAskTime).String())
		}
	}()

	for {
		time.Sleep(15 * time.Second)
		if uc.clearAll {
			order := agent.clearUnfinishedOrder()
			if order == nil {
				continue
			}
			agent.waitingOrder = order
			uc.placeOrder(agent, order)
		}
	}
}

func (uc *StreamUseCase) placeOrder(agent *TradeAgent, order *entity.Order) {
	if order.Price == 0 {
		log.Errorf("%s Order price is 0", order.StockNum)
		return
	}

	agent.waitingOrder = order

	// if out of trade in time, return
	if !uc.tradeInSwitch && (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) {
		// avoid stuck in the market
		agent.waitingOrder = nil
		return
	}

	bus.PublishTopicEvent(topicPlaceOrder, order)
	go agent.checkPlaceOrderStatus(order)
}

func (uc *StreamUseCase) checkTradeSwitch() {
	openTime := uc.basic.OpenTime
	endTime := uc.basic.EndTime
	tradeInEndTime := uc.basic.TradeInEndTime
	tradeOutEndTime := uc.basic.TradeOutEndTime

	for range time.NewTicker(5 * time.Second).C {
		now := time.Now()
		switch {
		case now.Before(openTime) || now.After(endTime) || (now.After(tradeInEndTime) && now.Before(tradeOutEndTime)):
			uc.tradeInSwitch = false
			uc.clearAll = false
		case now.After(openTime) && now.Before(tradeInEndTime):
			uc.tradeInSwitch = true
			uc.clearAll = false
		case now.After(tradeOutEndTime) && now.Before(endTime):
			uc.tradeInSwitch = false
			uc.clearAll = true
		}
	}
}

// GetTSESnapshot -.
func (uc *StreamUseCase) GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
	body, err := uc.grpcapi.GetStockSnapshotTSE()
	if err != nil {
		return nil, err
	}
	return &entity.StockSnapShot{
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
	}, nil
}

func (uc *StreamUseCase) realTimeAddTargets() error {
	data, err := uc.grpcapi.GetAllStockSnapshot()
	if err != nil {
		return err
	}

	// at least 200 snapshot to rank volume
	if len(data) < 200 {
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetTotalVolume() > data[j].GetTotalVolume()
	})
	data = data[:uc.targetCond.RealTimeRank]

	currentTargets := cc.GetTargets()
	targetsMap := make(map[string]*entity.Target)
	for _, t := range currentTargets {
		targetsMap[t.StockNum] = t
	}

	var newTargets []*entity.Target
	for _, c := range uc.targetCond.PriceVolumeLimit {
		for i, d := range data {
			stock := cc.GetStockDetail(d.GetCode())
			if stock == nil {
				continue
			}

			if !blackStockFilter(stock.Number, uc.targetCond) ||
				!blackCatagoryFilter(stock.Category, uc.targetCond) ||
				!targetFilter(d.GetClose(), d.GetTotalVolume(), c, true) {
				continue
			}

			if targetsMap[d.GetCode()] == nil {
				newTargets = append(newTargets, &entity.Target{
					Rank:     100 + i + 1,
					StockNum: d.GetCode(),
					Volume:   d.GetTotalVolume(),
					RealTime: true,
					TradeDay: uc.basic.TradeDay,
					Stock:    stock,
				})
			}
		}
	}

	if len(newTargets) != 0 {
		cc.AppendTargets(newTargets)
		bus.PublishTopicEvent(topicRealTimeTargets, newTargets, true)
	}
	return nil
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
			StockNum:        stockNum,
			StockName:       cc.GetStockDetail(stockNum).Name,
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
		})
	}

	for _, v := range stockNotExist {
		result = append(result, &entity.StockSnapShot{
			StockNum: v,
		})
	}

	return result, nil
}
