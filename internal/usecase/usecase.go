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
