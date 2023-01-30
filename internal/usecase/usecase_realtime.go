package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/target"
	"tmt/internal/usecase/module/trader"
	"tmt/internal/usecase/rabbit"
	"tmt/internal/usecase/topic"

	"github.com/google/uuid"
)

// RealTimeUseCase -.
type RealTimeUseCase struct {
	repo         RealTimeRepo
	grpcapi      RealTimegRPCAPI
	subgRPCAPI   SubscribegRPCAPI
	commonRabbit Rabbit
	futureRabbit Rabbit

	cfg          *config.Config
	targetFilter *target.Filter
	tradeIndex   *entity.TradeIndex

	mainFutureCode string
}

func NewRealTime(r RealTimeRepo, g RealTimegRPCAPI, s SubscribegRPCAPI) RealTime {
	cfg := config.GetConfig()
	uc := &RealTimeUseCase{
		repo:         r,
		commonRabbit: rabbit.NewRabbit(),
		futureRabbit: rabbit.NewRabbit(),
		grpcapi:      g,
		subgRPCAPI:   s,
		cfg:          cfg,
		targetFilter: target.NewFilter(cfg.TargetCond),
	}

	// unsubscriba all first
	if e := uc.UnSubscribeAll(); e != nil {
		logger.Fatal(e)
	}

	basic := cc.GetBasicInfo()
	uc.commonRabbit.FillAllBasic(basic.AllStocks, basic.AllFutures)
	uc.periodUpdateTradeIndex()

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	bus.SubscribeTopic(topic.TopicSubscribeStockTickTargets, uc.ReceiveStockSubscribeData, uc.SubscribeStockTick)
	bus.SubscribeTopic(topic.TopicUnSubscribeStockTickTargets, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)
	bus.SubscribeTopic(topic.TopicSubscribeFutureTickTargets, uc.ReceiveFutureSubscribeData, uc.SubscribeFutureTick)

	return uc
}

func (uc *RealTimeUseCase) GetTradeIndex() *entity.TradeIndex {
	return uc.tradeIndex
}

func (uc *RealTimeUseCase) periodUpdateTradeIndex() {
	uc.tradeIndex = &entity.TradeIndex{
		TSE:    entity.NewIndexStatus(),
		OTC:    entity.NewIndexStatus(),
		Nasdaq: entity.NewIndexStatus(),
		NF:     entity.NewIndexStatus(),
	}

	go uc.updateNasdaqIndex()
	go uc.updateNFIndex()
	go uc.updateTSEIndex()
	go uc.updateOTCIndex()
}

func (uc *RealTimeUseCase) updateNasdaqIndex() {
	for range time.NewTicker(time.Second * 5).C {
		if data, err := uc.GetNasdaqClose(); err != nil && !errors.Is(err, errNasdaqPriceAbnormal) {
			logger.Error(err)
		} else if data != nil {
			uc.tradeIndex.Nasdaq.UpdateIndexStatus(data.Price - data.Last)
		}
	}
}

func (uc *RealTimeUseCase) updateNFIndex() {
	for range time.NewTicker(time.Second * 5).C {
		if data, err := uc.GetNasdaqFutureClose(); err != nil && !errors.Is(err, errNFQPriceAbnormal) {
			logger.Error(err)
		} else if data != nil {
			uc.tradeIndex.NF.UpdateIndexStatus(data.Price - data.Last)
		}
	}
}

func (uc *RealTimeUseCase) updateTSEIndex() {
	for range time.NewTicker(time.Second * 3).C {
		if data, err := uc.GetTSESnapshot(context.Background()); err != nil {
			logger.Error(err)
		} else {
			uc.tradeIndex.TSE.UpdateIndexStatus(data.PriceChg)
		}
	}
}

func (uc *RealTimeUseCase) updateOTCIndex() {
	for range time.NewTicker(time.Second * 3).C {
		if data, err := uc.GetOTCSnapshot(context.Background()); err != nil {
			logger.Error(err)
		} else {
			uc.tradeIndex.OTC.UpdateIndexStatus(data.PriceChg)
		}
	}
}

// ReceiveEvent -.
func (uc *RealTimeUseCase) ReceiveEvent(ctx context.Context) {
	eventChan := make(chan *entity.SinopacEvent)
	go func() {
		for {
			event := <-eventChan
			if err := uc.repo.InsertEvent(ctx, event); err != nil {
				logger.Error(err)
			}

			if event.EventCode != 16 {
				logger.Warnf("EventCode: %d, Event: %s, ResoCode: %d, Info: %s", event.EventCode, event.Event, event.Response, event.Info)
			}
		}
	}()
	uc.commonRabbit.EventConsumer(eventChan)
}

// ReceiveOrderStatus -.
func (uc *RealTimeUseCase) ReceiveOrderStatus(ctx context.Context) {
	orderStatusChan := make(chan interface{})
	go func() {
		for {
			order := <-orderStatusChan
			switch t := order.(type) {
			case *entity.StockOrder:
				if cc.GetOrderByOrderID(t.OrderID) == nil {
					cc.SetOrderByOrderID(t.ToManual())
				}
				bus.PublishTopicEvent(topic.TopicInsertOrUpdateStockOrder, t)
			case *entity.FutureOrder:
				if cc.GetFutureOrderByOrderID(t.OrderID) == nil {
					cc.SetFutureOrderByOrderID(t.ToManual())
				}
				bus.PublishTopicEvent(topic.TopicInsertOrUpdateFutureOrder, t)
			}
		}
	}()
	uc.commonRabbit.AddOrderStatusChan(orderStatusChan, uuid.New().String())
	go uc.commonRabbit.OrderStatusConsumer()
	go uc.commonRabbit.OrderStatusArrConsumer()
}

// GetTSESnapshot -.
func (uc *RealTimeUseCase) GetTSESnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
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
func (uc *RealTimeUseCase) GetOTCSnapshot(ctx context.Context) (*entity.StockSnapShot, error) {
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

var (
	errNasdaqPriceAbnormal error = errors.New("nasdaq price abnormal")
	errNFQPriceAbnormal    error = errors.New("nfq price abnormal")
)

func (uc *RealTimeUseCase) GetNasdaqClose() (*entity.YahooPrice, error) {
	d, err := uc.grpcapi.GetNasdaq()
	if err != nil {
		return nil, err
	}

	if d.GetLast() == 0 || d.GetPrice() == 0 {
		return nil, errNasdaqPriceAbnormal
	}

	return &entity.YahooPrice{
		Last:      d.GetLast(),
		Price:     d.GetPrice(),
		UpdatedAt: time.Now(),
	}, nil
}

func (uc *RealTimeUseCase) GetNasdaqFutureClose() (*entity.YahooPrice, error) {
	d, err := uc.grpcapi.GetNasdaqFuture()
	if err != nil {
		return nil, err
	}

	if d.GetLast() == 0 || d.GetPrice() == 0 {
		return nil, errNFQPriceAbnormal
	}

	return &entity.YahooPrice{
		Last:      d.GetLast(),
		Price:     d.GetPrice(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetStockSnapshotByNumArr -.
func (uc *RealTimeUseCase) GetStockSnapshotByNumArr(stockNumArr []string) ([]*entity.StockSnapShot, error) {
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
func (uc *RealTimeUseCase) GetFutureSnapshotByCode(code string) (*entity.FutureSnapShot, error) {
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

// GetMainFuture -.
func (uc *RealTimeUseCase) GetMainFuture() *entity.Future {
	return cc.GetFutureDetail(uc.mainFutureCode)
}

// ReceiveStockSubscribeData - receive target data, start goroutine to trade
func (uc *RealTimeUseCase) ReceiveStockSubscribeData(targetArr []*entity.StockTarget) {
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
			r := rabbit.NewRabbit()
			go r.StockTickConsumer(agent.GetStockNum(), agent.GetTickChan())
			// go r.StockBidAskConsumer(agent.GetStockNum(), agent.GetBidAskChan())
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
			agent := trader.NewStockTrader(target.StockNum, uc.cfg.StockTradeSwitch, uc.cfg.StockAnalyze)
			agentChan <- agent
		}()
	}

	wg.Wait()
	close(agentChan)
	logger.Info("Stock trade room all start")
}

// ReceiveFutureSubscribeData -.
func (uc *RealTimeUseCase) ReceiveFutureSubscribeData(code string) {
	uc.mainFutureCode = code
	agent := trader.NewFutureTrader(code, uc.cfg.FutureTradeSwitch, uc.cfg.FutureAnalyze)

	go agent.TradingRoom()

	go uc.futureRabbit.FutureTickConsumer(code, agent.GetTickChan())
	// r := rabbit.NewRabbit()
	// go r.FutureTickConsumer(code, agent.GetTickChan())
	// go r.FutureBidAskConsumer(code, agent.GetBidAskChan())
	logger.Info("Future trade room start")
}

// NewFutureRealTimeConnection -.
func (uc *RealTimeUseCase) NewFutureRealTimeConnection(tickChan chan *entity.RealTimeFutureTick, connectionID string) {
	uc.futureRabbit.AddFutureTickChan(tickChan, connectionID)
}

// DeleteFutureRealTimeConnection -.
func (uc *RealTimeUseCase) DeleteFutureRealTimeConnection(connectionID string) {
	uc.futureRabbit.RemoveFutureTickChan(connectionID)
}

// NewOrderStatusConnection -.
func (uc *RealTimeUseCase) NewOrderStatusConnection(orderStatusChan chan interface{}, connectionID string) {
	uc.futureRabbit.AddOrderStatusChan(orderStatusChan, connectionID)
}

func (uc *RealTimeUseCase) DeleteOrderStatusConnection(connectionID string) {
	uc.futureRabbit.RemoveOrderStatusChan(connectionID)
}

// UnSubscribeAll -.
func (uc *RealTimeUseCase) UnSubscribeAll() error {
	result, err := uc.subgRPCAPI.UnSubscribeAllTick()
	if err != nil {
		return err
	}

	if m := result.GetErr(); m != "" {
		return errors.New(m)
	}

	result, err = uc.subgRPCAPI.UnSubscribeAllBidAsk()
	if err != nil {
		return err
	}

	if m := result.GetErr(); m != "" {
		return errors.New(m)
	}

	return nil
}

// SubscribeStockTick -.
func (uc *RealTimeUseCase) SubscribeStockTick(targetArr []*entity.StockTarget) error {
	if !uc.cfg.StockTradeSwitch.Subscribe {
		return nil
	}

	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
	}

	failSubNumArr, err := uc.subgRPCAPI.SubscribeStockTick(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// SubscribeStockBidAsk -.
func (uc *RealTimeUseCase) SubscribeStockBidAsk(targetArr []*entity.StockTarget) error {
	if !uc.cfg.StockTradeSwitch.Subscribe {
		return nil
	}

	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
	}

	failSubNumArr, err := uc.subgRPCAPI.SubscribeStockBidAsk(subArr)
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe fail %v", failSubNumArr)
	}

	return nil
}

// UnSubscribeStockTick -.
func (uc *RealTimeUseCase) UnSubscribeStockTick(stockNum string) error {
	failUnSubNumArr, err := uc.subgRPCAPI.UnSubscribeStockTick([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}

// UnSubscribeStockBidAsk -.
func (uc *RealTimeUseCase) UnSubscribeStockBidAsk(stockNum string) error {
	failUnSubNumArr, err := uc.subgRPCAPI.UnSubscribeStockBidAsk([]string{stockNum})
	if err != nil {
		return err
	}

	if len(failUnSubNumArr) != 0 {
		return fmt.Errorf("unsubscribe fail %v", failUnSubNumArr)
	}

	return nil
}

// SubscribeFutureTick -.
func (uc *RealTimeUseCase) SubscribeFutureTick(code string) error {
	if !uc.cfg.FutureTradeSwitch.Subscribe {
		return nil
	}

	failSubNumArr, err := uc.subgRPCAPI.SubscribeFutureTick([]string{code})
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe future fail %v", failSubNumArr)
	}

	return nil
}

// SubscribeFutureBidAsk -.
func (uc *RealTimeUseCase) SubscribeFutureBidAsk(code string) error {
	failSubNumArr, err := uc.subgRPCAPI.SubscribeFutureBidAsk([]string{code})
	if err != nil {
		return err
	}

	if len(failSubNumArr) != 0 {
		return fmt.Errorf("subscribe future fail %v", failSubNumArr)
	}

	return nil
}
