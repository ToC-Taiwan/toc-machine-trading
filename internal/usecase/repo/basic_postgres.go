// Package repo package repo
package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
)

// const _defaultEntityCap = 64

// BasicRepo -.
type BasicRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *BasicRepo {
	return &BasicRepo{pg}
}

// StoreStockDetail -.
func (r *BasicRepo) StoreStockDetail(ctx context.Context, t []*entity.Stock) error {
	builder := r.Builder.Insert("basic_stock").Columns("number, name, exchange, category, day_trade, last_close")
	for _, v := range t {
		builder = builder.Values(v.Number, v.Name, v.Exchange, v.Category, v.DayTrade, v.LastClose)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetAllStockDetail -.
func (r *BasicRepo) GetAllStockDetail(ctx context.Context) ([]*entity.Stock, error) {
	// builder := r.Builder.Insert("basic_stock").Columns("number, name, exchange, category, day_trade, last_close")
	// for _, v := range t {
	// 	builder = builder.Values(v.Number, v.Name, v.Exchange, v.Category, v.DayTrade, v.LastClose)
	// }

	// sql, args, err := builder.ToSql()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.Pool.Exec(ctx, sql, args...)
	// if err != nil {
	// 	return err
	// }
	return nil, nil
}
