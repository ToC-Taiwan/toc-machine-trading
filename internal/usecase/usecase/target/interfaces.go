package target

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=target_test

type Target interface {
	Init(logger *log.Log, cc *cache.Cache, bus *eventbus.Bus) Target
	GetTargets(ctx context.Context) []*entity.StockTarget
}

type TargetRepo interface {
	InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error
	QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error)
	QueryAllMXFFuture(ctx context.Context) ([]*entity.Future, error)
}
