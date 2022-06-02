package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
)

// StockUseCase -.
type StockUseCase struct {
	repo    StockRepo
	gRPCAPI StockgRPCAPI
}

// New -.
func New(r *repo.StockRepo, t *grpcapi.StockgRPCAPI) *StockUseCase {
	return &StockUseCase{
		repo:    r,
		gRPCAPI: t,
	}
}

// GetAllStockDetail -.
func (uc *StockUseCase) GetAllStockDetail(ctx context.Context) ([]*entity.Stock, error) {
	stockArr, err := uc.gRPCAPI.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	var stockDetail []*entity.Stock
	for _, v := range stockArr {
		stockDetail = append(stockDetail, v.ToStockEntity())
	}

	err = uc.repo.Store(context.Background(), stockDetail)
	if err != nil {
		return []*entity.Stock{}, err
	}

	return stockDetail, nil
}
