package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"tmt/internal/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpc"
	"tmt/internal/usecase/modules/cache"
	"tmt/internal/usecase/modules/calendar"
	"tmt/internal/usecase/modules/dt"
	"tmt/internal/usecase/modules/hadger"
	"tmt/internal/usecase/modules/quota"
	"tmt/internal/usecase/mqtt"
	"tmt/internal/usecase/repo"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

// RealTimeUseCase -.
type RealTimeUseCase struct {
	repo RealTimeRepo

	grpcapi    RealTimegRPCAPI
	subgRPCAPI SubscribegRPCAPI
	sc         TradegRPCAPI
	fg         TradegRPCAPI

	commonRabbit Rabbit
	futureRabbit Rabbit

	clientRabbitMap     map[string]Rabbit
	clientRabbitMapLock sync.RWMutex

	cfg        *config.Config
	quota      *quota.Quota
	tradeIndex *entity.TradeIndex

	mainFutureCode string

	stockSwitchChanMap      map[string]chan bool
	stockSwitchChanMapLock  sync.RWMutex
	futureSwitchChanMap     map[string]chan bool
	futureSwitchChanMapLock sync.RWMutex

	inventoryIsNotEmpty bool

	logger *log.Log
	cc     *cache.Cache
	bus    *eventbus.Bus
}

func NewRealTime(logger *log.Log, cc *cache.Cache, bus *eventbus.Bus) RealTime {
	cfg := config.Get()
	uc := &RealTimeUseCase{
		quota: quota.NewQuota(cfg.Quota),
		repo:  repo.NewRealTime(cfg.GetPostgresPool()),

		commonRabbit: mqtt.NewRabbit(cfg.GetRabbitConn()),
		futureRabbit: mqtt.NewRabbit(cfg.GetRabbitConn()),

		grpcapi:    grpc.NewRealTime(cfg.GetSinopacPool()),
		subgRPCAPI: grpc.NewSubscribe(cfg.GetSinopacPool()),

		sc: grpc.NewTrade(cfg.GetSinopacPool(), cfg.Simulation),
		fg: grpc.NewTrade(cfg.GetFuglePool(), cfg.Simulation),

		cfg:                 cfg,
		futureSwitchChanMap: make(map[string]chan bool),
		stockSwitchChanMap:  make(map[string]chan bool),

		clientRabbitMap: make(map[string]Rabbit),
	}
	uc.logger = logger
	uc.cc = cc
	uc.bus = bus

	// unsubscriba all first
	if e := uc.UnSubscribeAll(); e != nil {
		uc.logger.Fatal(e)
	}
	uc.checkFutureInventory()

	uc.commonRabbit.FillAllBasic(uc.cc.GetAllStockDetail(), uc.cc.GetAllFutureDetail())
	uc.futureRabbit.FillAllBasic(uc.cc.GetAllStockDetail(), uc.cc.GetAllFutureDetail())
	uc.periodUpdateTradeIndex()

	uc.checkFutureTradeSwitch()
	uc.checkStockTradeSwitch()

	go uc.ReceiveEvent(context.Background())
	go uc.ReceiveOrderStatus(context.Background())

	uc.bus.SubscribeAsync(topicSubscribeStockTickTargets, true, uc.ReceiveStockSubscribeData)
	uc.bus.SubscribeAsync(topicUnSubscribeStockTickTargets, false, uc.UnSubscribeStockTick, uc.UnSubscribeStockBidAsk)
	uc.bus.SubscribeAsync(topicSubscribeFutureTickTargets, true, uc.SetMainFuture)

	return uc
}

func (uc *RealTimeUseCase) checkFutureInventory() {
	position, err := uc.sc.GetFuturePosition()
	if err != nil {
		uc.logger.Fatal(err)
		return
	}

	for _, v := range position.GetPositionArr() {
		if uc.cc.GetFutureDetail(v.GetCode()) != nil {
			uc.logger.Warnf("future inventory is not empty, code: %s", v.GetCode())
			uc.inventoryIsNotEmpty = true
			return
		}
	}
}

func (uc *RealTimeUseCase) checkStockTradeSwitch() {
	if !uc.cfg.TradeStock.AllowTrade || uc.cfg.ManualTrade {
		return
	}
	stockTradeDay := calendar.Get().GetStockTradeDay().TradeDay
	openTime := stockTradeDay.Add(9 * time.Hour).Add(time.Duration(uc.cfg.TradeStock.HoldTimeFromOpen) * time.Second)
	tradeInEndTime := stockTradeDay.Add(9 * time.Hour).Add(time.Duration(uc.cfg.TradeStock.TradeInEndTime) * time.Minute)

	go func() {
		for range time.NewTicker(30 * time.Second).C {
			now := time.Now()
			var tempSwitch bool
			switch {
			case now.Before(openTime) || now.After(tradeInEndTime):
				tempSwitch = false
			case now.After(openTime) && now.Before(tradeInEndTime):
				tempSwitch = true
			}

			uc.stockSwitchChanMapLock.RLock()
			for _, ch := range uc.stockSwitchChanMap {
				ch <- tempSwitch
			}
			uc.stockSwitchChanMapLock.RUnlock()
		}
	}()
}

func (uc *RealTimeUseCase) checkFutureTradeSwitch() {
	if !uc.cfg.TradeFuture.AllowTrade || uc.inventoryIsNotEmpty || uc.cfg.ManualTrade {
		return
	}
	futureTradeDay := calendar.Get().GetFutureTradeDay()

	timeRange := [][]time.Time{}
	firstStart := futureTradeDay.StartTime
	secondStart := futureTradeDay.EndTime.Add(-300 * time.Minute)

	timeRange = append(timeRange, []time.Time{firstStart, firstStart.Add(time.Duration(uc.cfg.TradeFuture.TradeTimeRange.FirstPartDuration) * time.Minute)})
	timeRange = append(timeRange, []time.Time{secondStart, secondStart.Add(time.Duration(uc.cfg.TradeFuture.TradeTimeRange.SecondPartDuration) * time.Minute)})

	go func() {
		for range time.NewTicker(30 * time.Second).C {
			now := time.Now()
			var tempSwitch bool
			for _, rangeTime := range timeRange {
				if now.After(rangeTime[0]) && now.Before(rangeTime[1]) {
					tempSwitch = true
				}
			}

			uc.futureSwitchChanMapLock.RLock()
			for _, ch := range uc.futureSwitchChanMap {
				ch <- tempSwitch
			}
			uc.futureSwitchChanMapLock.RUnlock()
		}
	}()
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
			uc.logger.Error(err)
		} else if data != nil {
			uc.tradeIndex.Nasdaq.UpdateIndexStatus(data.Price - data.Last)
		}
	}
}

func (uc *RealTimeUseCase) updateNFIndex() {
	for range time.NewTicker(time.Second * 5).C {
		if data, err := uc.GetNasdaqFutureClose(); err != nil && !errors.Is(err, errNFQPriceAbnormal) {
			uc.logger.Error(err)
		} else if data != nil {
			uc.tradeIndex.NF.UpdateIndexStatus(data.Price - data.Last)
		}
	}
}

func (uc *RealTimeUseCase) updateTSEIndex() {
	for range time.NewTicker(time.Second * 3).C {
		if data, err := uc.GetTSESnapshot(context.Background()); err != nil {
			uc.logger.Error(err)
		} else {
			uc.tradeIndex.TSE.UpdateIndexStatus(data.PriceChg)
		}
	}
}

func (uc *RealTimeUseCase) updateOTCIndex() {
	for range time.NewTicker(time.Second * 3).C {
		if data, err := uc.GetOTCSnapshot(context.Background()); err != nil {
			uc.logger.Error(err)
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
				uc.logger.Error(err)
			}

			if event.EventCode != 16 {
				uc.logger.Warnf("EventCode: %d, Event: %s, ResoCode: %d, Info: %s", event.EventCode, event.Event, event.Response, event.Info)
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
				uc.bus.PublishTopicEvent(topicInsertOrUpdateStockOrder, t)
			case *entity.FutureOrder:
				uc.bus.PublishTopicEvent(topicInsertOrUpdateFutureOrder, t)
			}
		}
	}()
	go uc.commonRabbit.OrderStatusConsumer(orderStatusChan)
	go uc.commonRabbit.OrderStatusArrConsumer(orderStatusChan)
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
		if uc.cc.GetStockDetail(s) == nil {
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
			StockName: uc.cc.GetStockDetail(stockNum).Name,
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
		FutureName: uc.cc.GetFutureDetail(code).Name,
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
	return uc.cc.GetFutureDetail(uc.mainFutureCode)
}

// ReceiveStockSubscribeData - receive target data, start goroutine to trade
func (uc *RealTimeUseCase) ReceiveStockSubscribeData(targetArr []*entity.StockTarget) {
	defer uc.SubscribeStockTick(targetArr)
	notifyChanMap := make(map[string]chan *entity.StockOrder)
	for _, t := range targetArr {
		hadger := hadger.NewHadgerStock(
			t.StockNum,
			uc.sc.(*grpc.TradegRPCAPI),
			uc.fg.(*grpc.TradegRPCAPI),
			uc.quota,
			&uc.cfg.TradeStock,
		)
		notifyChanMap[t.StockNum] = hadger.Notify()

		uc.stockSwitchChanMapLock.Lock()
		uc.stockSwitchChanMap[t.StockNum] = hadger.SwitchChan()
		uc.stockSwitchChanMapLock.Unlock()

		uc.logger.Infof("Stock room %s <-> %s <-> %s", t.Stock.Name, t.Stock.Future.Name, t.Stock.Future.Code)
		r := mqtt.NewRabbit(uc.cfg.GetRabbitConn())
		go r.StockTickConsumer(t.StockNum, hadger.TickChan())
	}

	orderStatusChan := make(chan interface{})
	go func() {
		for {
			order := <-orderStatusChan
			o, ok := order.(*entity.StockOrder)
			if !ok {
				continue
			}

			if ch := notifyChanMap[o.StockNum]; ch != nil {
				ch <- o
			}
		}
	}()
	hr := mqtt.NewRabbit(uc.cfg.GetRabbitConn())
	hr.FillAllBasic(uc.cc.GetAllStockDetail(), uc.cc.GetAllFutureDetail())
	go hr.OrderStatusConsumer(orderStatusChan)
	go hr.OrderStatusArrConsumer(orderStatusChan)
}

// ReceiveFutureSubscribeData -.
func (uc *RealTimeUseCase) ReceiveFutureSubscribeData(code string) {
	defer uc.SubscribeFutureTick(code)
	t := dt.NewDTFuture(
		code,
		uc.sc.(*grpc.TradegRPCAPI),
		&uc.cfg.TradeFuture,
	)

	uc.futureSwitchChanMapLock.Lock()
	uc.futureSwitchChanMap[code] = t.SwitchChan()
	uc.futureSwitchChanMapLock.Unlock()

	orderStatusChan := make(chan interface{})
	go uc.futureRabbit.FutureTickConsumer(code, t.TickChan())
	go func() {
		ch := t.Notify()
		for {
			order := <-orderStatusChan
			o, ok := order.(*entity.FutureOrder)
			if !ok || o.Code != code {
				continue
			}
			ch <- o
		}
	}()
	go uc.futureRabbit.OrderStatusConsumer(orderStatusChan)
	go uc.futureRabbit.OrderStatusArrConsumer(orderStatusChan)
}

func (uc *RealTimeUseCase) SetMainFuture(code string) {
	defer uc.ReceiveFutureSubscribeData(code)

	if uc.mainFutureCode != "" {
		uc.logger.Fatal("main future code already set, can't set again")
	}

	uc.mainFutureCode = code
}

func (uc *RealTimeUseCase) NewFutureRealTimeClient(tickChan chan *entity.RealTimeFutureTick, orderStatusChan chan interface{}, connectionID string) {
	r := mqtt.NewRabbit(uc.cfg.GetRabbitConn())
	r.FillAllBasic(uc.cc.GetAllStockDetail(), uc.cc.GetAllFutureDetail())

	uc.clientRabbitMapLock.Lock()
	uc.clientRabbitMap[connectionID] = r
	uc.clientRabbitMapLock.Unlock()

	go r.OrderStatusConsumer(orderStatusChan)
	go r.OrderStatusArrConsumer(orderStatusChan)
	go r.FutureTickConsumer(uc.mainFutureCode, tickChan)
}

func (uc *RealTimeUseCase) DeleteFutureRealTimeClient(connectionID string) {
	uc.clientRabbitMapLock.Lock()
	defer uc.clientRabbitMapLock.Unlock()
	if r, ok := uc.clientRabbitMap[connectionID]; ok {
		r.Close()
		delete(uc.clientRabbitMap, connectionID)
	}
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
func (uc *RealTimeUseCase) SubscribeStockTick(targetArr []*entity.StockTarget) {
	if !uc.cfg.TradeStock.Subscribe {
		return
	}

	var subArr []string
	for _, v := range targetArr {
		subArr = append(subArr, v.StockNum)
	}

	failSubNumArr, err := uc.subgRPCAPI.SubscribeStockTick(subArr, uc.cfg.TradeStock.Odd)
	if err != nil {
		uc.logger.Error(err)
		return
	}

	if len(failSubNumArr) != 0 {
		uc.logger.Errorf("subscribe fail %v", failSubNumArr)
	}
}

// SubscribeStockBidAsk -.
func (uc *RealTimeUseCase) SubscribeStockBidAsk(targetArr []*entity.StockTarget) error {
	if !uc.cfg.TradeStock.Subscribe {
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
func (uc *RealTimeUseCase) SubscribeFutureTick(code string) {
	if !uc.cfg.TradeFuture.Subscribe {
		return
	}

	failSubNumArr, err := uc.subgRPCAPI.SubscribeFutureTick([]string{code})
	if err != nil {
		uc.logger.Error(err)
		return
	}

	if len(failSubNumArr) != 0 {
		uc.logger.Errorf("subscribe future fail %v", failSubNumArr)
	}
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
