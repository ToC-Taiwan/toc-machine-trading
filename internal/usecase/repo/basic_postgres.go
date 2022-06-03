// Package repo package repo
package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
)

const _defaultEntityCap = 1800

// BasicRepo -.
type BasicRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *BasicRepo {
	return &BasicRepo{pg}
}

// InsertStock -.
func (r *BasicRepo) InsertStock(ctx context.Context, t []*entity.Stock) error {
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

// QueryAllStock -.
func (r *BasicRepo) QueryAllStock(ctx context.Context) ([]*entity.Stock, error) {
	sql, _, err := r.Builder.
		Select("*").
		From("basic_stock").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.Stock, 0, _defaultEntityCap)
	for rows.Next() {
		e := &entity.Stock{}
		err = rows.Scan(
			&e.ID,
			&e.Number,
			&e.Name,
			&e.Exchange,
			&e.Category,
			&e.DayTrade,
			&e.LastClose,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}
