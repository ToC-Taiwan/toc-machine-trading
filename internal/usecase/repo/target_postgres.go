package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
)

const (
	tableNameTarget string = "basic_targets"
)

// TargetRepo -.
type TargetRepo struct {
	*postgres.Postgres
}

// NewTarget -.
func NewTarget(pg *postgres.Postgres) *TargetRepo {
	return &TargetRepo{pg}
}

// InsertTargetArr -.
func (r *TargetRepo) InsertTargetArr(ctx context.Context, t []*entity.Target) error {
	builder := r.Builder.Insert(tableNameTarget).Columns("stock_num, trade_day, rank, volume, subscribe, real_time_add")
	for _, v := range t {
		builder = builder.Values(v.StockNum, v.TradeDay, v.Rank, v.Volume, v.Subscribe, v.RealTimeAdd)
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

// QueryAllTargetByTradeDay -.
func (r *TargetRepo) QueryAllTargetByTradeDay(ctx context.Context) ([]*entity.Target, error) {
	sql, _, err := r.Builder.
		Select("*").From(tableNameTarget).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.Target, 0, 256)
	for rows.Next() {
		e := &entity.Target{}
		err = rows.Scan(
			&e.StockNum, &e.TradeDay, &e.Rank, &e.Volume, &e.Subscribe, &e.RealTimeAdd,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}
