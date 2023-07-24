package basic

import (
	"context"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/pb"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=basic_test

type Basic interface {
	Init(logger *log.Log, cc *cache.Cache) Basic
	GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error)
	GetConfig() *config.Config
	GetShioajiUsage() (*entity.ShioajiUsage, error)
}

type BasicRepo interface {
	QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error)
	InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
	UpdateAllStockDayTradeToNo(ctx context.Context) error

	QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error)
	InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error

	QueryAllFuture(ctx context.Context) (map[string]*entity.Future, error)
	QueryFutureByLikeName(ctx context.Context, name string) ([]*entity.Future, error)
	InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error

	QueryAllOption(ctx context.Context) (map[string]*entity.Option, error)
	QueryOptionByLikeName(ctx context.Context, name string) ([]*entity.Option, error)
	InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error
}

type BasicgRPCAPI interface {
	CreateLongConnection() error
	Terminate() error
	Login() error
	CheckUsage() (*pb.ShioajiUsage, error)
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	GetAllFutureDetail() ([]*pb.FutureDetailMessage, error)
	GetAllOptionDetail() ([]*pb.OptionDetailMessage, error)
}
