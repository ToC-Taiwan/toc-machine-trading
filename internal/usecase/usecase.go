package usecase

import (
	"fmt"

	"tmt/cmd/config"
	"tmt/internal/usecase/cache"
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
	sc  *grpc.ConnPool
	fg  *grpc.ConnPool
	cfg *config.Config
}

func NewUseCaseBase(cfg *config.Config) *UseCaseBase {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Database.URL, cfg.Database.DBName),
		postgres.MaxPoolSize(cfg.Database.PoolMax),
		postgres.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
		grpc.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		cfg.Fugle.URL,
		grpc.MaxPoolSize(cfg.Fugle.PoolMax),
		grpc.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	uc := &UseCaseBase{
		pg:  pg,
		sc:  sc,
		fg:  fg,
		cfg: cfg,
	}

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	return uc
}

func (u *UseCaseBase) Close() {
	u.pg.Close()
	logger.Warn("TMT is shutting down")
}
