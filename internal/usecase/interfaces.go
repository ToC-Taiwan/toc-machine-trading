// Package usecase package usecase
package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/pb"
)

// Stock -.
type Stock interface {
	GetAllStockDetail(ctx context.Context) ([]*entity.Stock, error)
}

// StockRepo -.
type StockRepo interface {
	Store(ctx context.Context, t []*entity.Stock) error
}

// StockgRPCAPI -.
type StockgRPCAPI interface {
	GetAllStockDetail() ([]*pb.StockDetailMessage, error)
}
