package usecase

import (
	"tmt/cmd/config"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/rabbit"
	"tmt/pkg/grpc"

	"github.com/robfig/cron/v3"
)

type Health struct {
	rabbit Rabbit
}

func StartHealthCheck() *Health {
	cfg := config.Get()
	uc := &Health{
		rabbit: rabbit.NewRabbit(cfg.GetRabbitConn()),
	}

	if e := uc.setupCronJob(); e != nil {
		logger.Fatal(e)
	}

	go uc.healthCheckforSinopac(cfg.GetSinopacPool(), cfg.Development)
	go uc.healthCheckforFugle(cfg.GetFuglePool(), cfg.Development)

	logger.Warn("TMT is running")
	logger.Warnf("Simulation Mode: %v", cfg.Simulation)

	return uc
}

func (u *Health) setupCronJob() error {
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

func (u *Health) healthCheckforSinopac(client *grpc.ConnPool, devMode bool) {
	sc := grpcapi.NewBasic(client, devMode)
	err := sc.Heartbeat()
	if err != nil {
		u.rabbit.PublishTerminate()
		logger.Fatal("sinopac healthcheck fail, publish terminate")
	}
}

func (u *Health) healthCheckforFugle(client *grpc.ConnPool, devMode bool) {
	fg := grpcapi.NewBasic(client, devMode)
	err := fg.Heartbeat()
	if err != nil {
		u.rabbit.PublishTerminate()
		logger.Fatal("fugle healthcheck fail, publish terminate")
	}
}

func (u *Health) Stop() {
	u.rabbit.PublishTerminate()
}
