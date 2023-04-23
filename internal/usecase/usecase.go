package usecase

import (
	"fmt"

	"tmt/cmd/config"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/rabbit"
	"tmt/pkg/eventbus"
	"tmt/pkg/grpc"
	"tmt/pkg/log"
	"tmt/pkg/postgres"

	"github.com/robfig/cron/v3"
)

var (
	logger = log.Get()
	cc     = cache.Get()
	bus    = eventbus.Get()
)

type UseCaseBase struct {
	pg     *postgres.Postgres
	sc     *grpc.ConnPool
	fg     *grpc.ConnPool
	cfg    *config.Config
	rabbit Rabbit
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
		pg:     pg,
		sc:     sc,
		fg:     fg,
		cfg:    cfg,
		rabbit: rabbit.NewRabbit(cfg.RabbitMQ),
	}

	if e := uc.SetupCronJob(); e != nil {
		logger.Fatal(e)
	}
	go uc.healthCheckforSinopac()
	go uc.healthCheckforFugle()

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	return uc
}

func (u *UseCaseBase) Close() {
	u.pg.Close()
	logger.Warn("TMT is shutting down")
}

func (u *UseCaseBase) SetupCronJob() error {
	c := cron.New()
	if _, e := c.AddFunc("20 8 * * *", u.rabbit.PublishTerminate); e != nil {
		return e
	}
	if _, e := c.AddFunc("40 14 * * *", u.rabbit.PublishTerminate); e != nil {
		return e
	}
	c.Start()
	return nil
}

func (u *UseCaseBase) healthCheckforSinopac() {
	sc := grpcapi.NewBasic(u.sc, u.cfg.Development)
	err := sc.Heartbeat()
	if err != nil {
		u.rabbit.PublishTerminate()
		logger.Fatal("sinopac healthcheck fail, publish terminate")
	}
}

func (u *UseCaseBase) healthCheckforFugle() {
	fg := grpcapi.NewBasic(u.fg, u.cfg.Development)
	err := fg.Heartbeat()
	if err != nil {
		u.rabbit.PublishTerminate()
		logger.Fatal("fugle healthcheck fail, publish terminate")
	}
}
