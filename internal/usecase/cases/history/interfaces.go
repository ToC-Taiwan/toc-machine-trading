package history

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pb"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=history_test

type History interface {
	Init(logger *log.Log, cc *usecase.Cache, bus *eventbus.Bus) History
	GetTradeDay() time.Time

	GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar
	FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error)
}

type HistoryRepo interface {
	InsertHistoryCloseArr(ctx context.Context, t []*entity.StockHistoryClose) error
	QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.StockHistoryClose, error)
	InsertHistoryTickArr(ctx context.Context, t []*entity.StockHistoryTick) error
	QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryTick, error)
	InsertHistoryKbarArr(ctx context.Context, t []*entity.StockHistoryKbar) error
	QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryKbar, error)
	InsertQuaterMA(ctx context.Context, t *entity.StockHistoryAnalyze) error
	QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.StockHistoryAnalyze, error)
	DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error
	InsertFutureHistoryTickArr(ctx context.Context, t []*entity.FutureHistoryTick) error
	QueryFutureHistoryTickArrByTime(ctx context.Context, code string, startTime, endTime time.Time) ([]*entity.FutureHistoryTick, error)
	InsertFutureHistoryClose(ctx context.Context, c *entity.FutureHistoryClose) error
	QueryFutureHistoryCloseByDate(ctx context.Context, code string, tradeDay time.Time) (*entity.FutureHistoryClose, error)
}

type HistorygRPCAPI interface {
	GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error)
	GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error)
	GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error)
	GetStockHistoryCloseByDateArr(stockNumArr []string, date []string) ([]*pb.HistoryCloseMessage, error)
	GetStockTSEHistoryTick(date string) ([]*pb.HistoryTickMessage, error)
	GetStockTSEHistoryKbar(date string) ([]*pb.HistoryKbarMessage, error)
	GetStockTSEHistoryClose(date string) ([]*pb.HistoryCloseMessage, error)
	GetFutureHistoryTick(codeArr []string, date string) ([]*pb.HistoryTickMessage, error)
	GetFutureHistoryKbar(codeArr []string, date string) ([]*pb.HistoryKbarMessage, error)
	GetFutureHistoryClose(codeArr []string, date string) ([]*pb.HistoryCloseMessage, error)
}
