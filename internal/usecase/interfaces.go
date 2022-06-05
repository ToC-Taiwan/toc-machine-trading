// Package usecase package usecase
package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/pb"
)

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
		InserOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	}
)

// type (
// 	// Stream -.
// 	Stream interface{}

// 	// StreamRepo -.
// 	StreamRepo interface{}

// 	// StreamgRPCAPI -.
// 	StreamgRPCAPI interface{}
// )
