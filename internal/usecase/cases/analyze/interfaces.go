package analyze

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=analyze_test

type Analyze interface {
	Init(logger *log.Log, cc *usecase.Cache, bus *eventbus.Bus) Analyze
	GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock
}
