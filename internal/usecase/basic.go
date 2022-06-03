package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo    BasicRepo
	gRPCAPI BasicgRPCAPI
	bus     *eventbus.Bus
}

// NewBasic -.
func NewBasic(r *repo.BasicRepo, t *grpcapi.BasicgRPCAPI, bus *eventbus.Bus) *BasicUseCase {
	return &BasicUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}
}

// GetAllSinopacStockAndUpdateRepo -.
func (uc *BasicUseCase) GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error) {
	stockArr, err := uc.gRPCAPI.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	var stockDetail []*entity.Stock
	for _, v := range stockArr {
		new := entity.Stock{}
		stockDetail = append(stockDetail, new.FromProto(v))
	}

	err = uc.repo.InsertStock(context.Background(), stockDetail)
	if err != nil {
		return []*entity.Stock{}, err
	}

	return stockDetail, nil
}

// GetAllRepoStock -.
func (uc *BasicUseCase) GetAllRepoStock(ctx context.Context) ([]*entity.Stock, error) {
	data, err := uc.repo.QueryAllStock(context.Background())
	if err != nil {
		return []*entity.Stock{}, err
	}
	return data, nil
}
