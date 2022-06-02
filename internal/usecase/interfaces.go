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
		GetAllStockDetail(ctx context.Context) ([]*entity.Stock, error)
	}

	// BasicRepo -.
	BasicRepo interface {
		GetAllStockDetail(ctx context.Context) ([]*entity.Stock, error)
		StoreStockDetail(ctx context.Context, t []*entity.Stock) error
	}

	// BasicgRPCAPI -.
	BasicgRPCAPI interface {
		GetAllStockDetail() ([]*pb.StockDetailMessage, error)
	}
)

type (
	// Stream -.
	Stream interface{}

	// StreamRepo -.
	StreamRepo interface{}

	// StreamgRPCAPI -.
	StreamgRPCAPI interface{}
)
