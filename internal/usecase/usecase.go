package usecase

import (
	"fmt"

	"tmt/cmd/config"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/slack"
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
	pg    *postgres.Postgres
	sc    *grpc.ConnPool
	fg    *grpc.ConnPool
	cfg   *config.Config
	slack *slack.Slack
}

func NewUseCaseBase(cfg *config.Config) *UseCaseBase {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", cfg.Database.URL, cfg.Database.DBName),
		postgres.MaxPoolSize(cfg.Database.PoolMax),
		postgres.Logger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		cfg.Sinopac.URL,
		grpc.MaxPoolSize(cfg.Sinopac.PoolMax),
		grpc.Logger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		cfg.Fugle.URL,
		grpc.MaxPoolSize(cfg.Fugle.PoolMax),
		grpc.Logger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	uc := &UseCaseBase{
		pg:    pg,
		sc:    sc,
		fg:    fg,
		cfg:   cfg,
		slack: slack.NewSlack(cfg.Slack.Token, cfg.Slack.ChannelID),
	}

	uc.slack.PostMessage("TMT is running :honeybee:")
	uc.slack.PostMessage(fmt.Sprintf("Simulation Mode: %v :computer:", cfg.Simulation))

	return uc
}

func (u *UseCaseBase) Close() {
	u.pg.Close()
	u.slack.PostMessage("TMT is shutting down")
}
