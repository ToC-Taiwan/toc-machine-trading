package usecase

import (
	"context"
	"fmt"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/target"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/rabbit"
	"tmt/internal/usecase/repo"
	"tmt/internal/usecase/topic"
	"tmt/pkg/eventbus"
	"tmt/pkg/grpc"
	"tmt/pkg/log"
	"tmt/pkg/postgres"
)

var (
	logger = log.Get()
	cc     = cache.Get()
	bus    = eventbus.Get()
)

type UseCaseBase struct {
	pg  *postgres.Postgres
	sc  *grpc.Connection
	fg  *grpc.Connection
	cfg *config.Config
}

func NewUseCaseBase(cfg *config.Config) *UseCaseBase {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Database.URL, cfg.Database.DBName),
		postgres.MaxPoolSize(cfg.Database.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		cfg.Fugle.URL,
		grpc.MaxPoolSize(cfg.Fugle.PoolMax),
	)
	if err != nil {
		logger.Fatal(err)
	}

	return &UseCaseBase{
		pg:  pg,
		sc:  sc,
		fg:  fg,
		cfg: cfg,
	}
}

func (u *UseCaseBase) Close() {
	u.pg.Close()
}

func (u *UseCaseBase) NewAnalyze() Analyze {
	uc := &AnalyzeUseCase{
		repo:             repo.NewHistory(u.pg),
		lastBelowMAStock: make(map[string]*entity.StockHistoryAnalyze),
		rebornMap:        make(map[time.Time][]entity.Stock),
		tradeDay:         tradeday.NewTradeDay(),
	}

	bus.SubscribeTopic(topic.TopicAnalyzeStockTargets, uc.findBelowQuaterMATargets)
	return uc
}

func (u *UseCaseBase) NewBasic() Basic {
	uc := &BasicUseCase{
		repo:     repo.NewBasic(u.pg),
		sc:       grpcapi.NewBasic(u.sc),
		fg:       grpcapi.NewBasic(u.fg),
		tradeDay: tradeday.NewTradeDay(),
		cfg:      u.cfg,
	}

	go uc.healthCheckforSinopac()
	go uc.healthCheckforFugle()

	if err := uc.importCalendarDate(context.Background()); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoStock(); err != nil {
		logger.Fatal(err)
	}

	if _, err := uc.updateRepoFuture(); err != nil {
		logger.Fatal(err)
	}

	uc.fillBasicInfo()

	return uc
}

// NewHistory -.
func (u *UseCaseBase) NewHistory() History {
	uc := &HistoryUseCase{
		repo:            repo.NewHistory(u.pg),
		grpcapi:         grpcapi.NewHistory(u.sc),
		fetchList:       make(map[string]*entity.StockTarget),
		tradeDay:        tradeday.NewTradeDay(),
		stockAnalyzeCfg: u.cfg.StockAnalyze,
	}

	uc.basic = cc.GetBasicInfo()
	bus.SubscribeTopic(topic.TopicFetchStockHistory, uc.FetchHistory)
	bus.SubscribeTopic(topic.TopicSubscribeFutureTickTargets, uc.updateMainFutureCode)
	return uc
}

func (u *UseCaseBase) NewRealTime() RealTime {
	cfg := u.cfg
	uc := &RealTimeUseCase{
		repo:         repo.NewRealTime(u.pg),
		commonRabbit: rabbit.NewRabbit(cfg.RabbitMQ),
		futureRabbit: rabbit.NewRabbit(cfg.RabbitMQ),
		grpcapi:      grpcapi.NewRealTime(u.sc),
		subgRPCAPI:   grpcapi.NewSubscribe(u.sc),
		cfg:          cfg,
		sc:           grpcapi.NewTrade(u.sc, cfg.Simulation),
		fg:           grpcapi.NewTrade(u.fg, cfg.Simulation),
		targetFilter: target.NewFilter(cfg.TargetCond),
		quota:        quota.NewQuota(cfg.Quota),
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

func (u *UseCaseBase) NewTarget() Target {
	cfg := u.cfg
	basic := cc.GetBasicInfo()
	uc := &TargetUseCase{
		repo:         repo.NewTarget(u.pg),
		gRPCAPI:      grpcapi.NewRealTime(u.sc),
		cfg:          cfg,
		basic:        basic,
		tradeDay:     tradeday.NewTradeDay(),
		targetFilter: target.NewFilter(cfg.TargetCond),
	}

	go uc.checkStockTradeSwitch()
	go uc.checkFutureTradeSwitch()

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

	go func() {
		time.Sleep(time.Until(basic.TradeDay.Add(time.Hour * 9)))
		for range time.NewTicker(time.Second * 60).C {
			if uc.stockTradeInSwitch {
				if err := uc.realTimeAddTargets(); err != nil {
					logger.Fatal(err)
				}
			}
		}
	}()

	return uc
}

func (u *UseCaseBase) NewTrade() Trade {
	cfg := u.cfg
	tradeDay := tradeday.NewTradeDay()

	uc := &TradeUseCase{
		simTrade: cfg.Simulation,

		sc:    grpcapi.NewTrade(u.sc, cfg.Simulation),
		fg:    grpcapi.NewTrade(u.fg, cfg.Simulation),
		repo:  repo.NewTrade(u.pg),
		quota: quota.NewQuota(cfg.Quota),

		tradeDay:       tradeDay,
		stockTradeDay:  tradeDay.GetStockTradeDay(),
		futureTradeDay: tradeDay.GetFutureTradeDay(),
	}

	bus.SubscribeTopic(topic.TopicInsertOrUpdateStockOrder, uc.updateStockOrderCacheAndInsertDB)
	bus.SubscribeTopic(topic.TopicInsertOrUpdateFutureOrder, uc.updateFutureOrderCacheAndInsertDB)

	if uc.simTrade {
		go uc.askSimulateOrderStatus()
	} else {
		go uc.askOrderStatus()
	}

	go uc.updateAllTradeBalance()
	return uc
}
