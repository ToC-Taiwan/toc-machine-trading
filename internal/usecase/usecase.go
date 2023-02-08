package usecase

import (
	"fmt"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/module/tradeday"
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
		tradeDay:         tradeday.Get(),
	}

	bus.SubscribeTopic(topic.TopicAnalyzeStockTargets, uc.findBelowQuaterMATargets)
	return uc
}
