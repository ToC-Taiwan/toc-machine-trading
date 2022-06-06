// Package usecase package usecase
package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test
type (
	// Basic -.
	Basic interface {
		GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error)
		GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error)
	}

	// BasicRepo -.
	BasicRepo interface {
		QueryAllStock(ctx context.Context) ([]*entity.Stock, error)
		InserOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error
		QueryAllTradeDay(ctx context.Context) ([]*entity.CalendarDate, error)
		InserOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	}
)

type (
	// Stream -.
	Stream interface {
		ReceiveEvent(ctx context.Context)
	}

	// StreamRepo -.
	StreamRepo interface {
		InsertEvent(ctx context.Context, t *entity.SinopacEvent) error
	}

	// StreamgRPCAPI -.
	StreamgRPCAPI interface {
		EventChannel(eventChan chan *entity.SinopacEvent) error
		BidAskChannel(bidAskChan chan *entity.RealTimeBidAsk) error
		TickChannel(tickChan chan *entity.RealTimeTick) error
	}
)
