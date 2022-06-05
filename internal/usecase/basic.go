package usecase

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/grpcapi"
	"toc-machine-trading/internal/usecase/repo"
	"toc-machine-trading/pkg/eventbus"
	"toc-machine-trading/pkg/logger"
)

// BasicUseCase -.
type BasicUseCase struct {
	repo    BasicRepo
	gRPCAPI BasicgRPCAPI
	bus     *eventbus.Bus
}

// NewBasic -.
func NewBasic(r *repo.BasicRepo, t *grpcapi.BasicgRPCAPI, bus *eventbus.Bus) *BasicUseCase {
	uc := &BasicUseCase{
		repo:    r,
		gRPCAPI: t,
		bus:     bus,
	}

	if _, err := uc.GetAllSinopacStockAndUpdateRepo(context.Background()); err != nil {
		logger.Get().Panic(err)
	}

	return uc
}

// GetAllSinopacStockAndUpdateRepo -.
func (uc *BasicUseCase) GetAllSinopacStockAndUpdateRepo(ctx context.Context) ([]*entity.Stock, error) {
	stockArr, err := uc.gRPCAPI.GetAllStockDetail()
	if err != nil {
		return []*entity.Stock{}, err
	}

	var stockDetail []*entity.Stock
	for _, v := range stockArr {
		stock := new(entity.Stock).FromProto(v)
		stockDetail = append(stockDetail, stock)

		// save to cache
		SetStockDetail(stock)
	}

	err = uc.repo.InserOrUpdatetStockArr(context.Background(), stockDetail)
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

	for _, s := range data {
		// save to cache
		SetStockDetail(s)
	}

	return data, nil
}
