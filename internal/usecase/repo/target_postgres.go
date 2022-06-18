package repo

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
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

	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

// QueryTargetsByTradeDay -.
func (r *TargetRepo) QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.Target, error) {
	sql, args, err := r.Builder.
		Select("id, rank, volume, subscribe, real_time_add, trade_day, stock_num, number, name, exchange, category, day_trade, last_close").
		From(tableNameTarget).
		Where("trade_day = ?", tradeDay).
		Join("basic_stock ON basic_targets.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.Target, 0, 256)
	for rows.Next() {
		e := entity.Target{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.ID, &e.Rank, &e.Volume, &e.Subscribe, &e.RealTimeAdd, &e.TradeDay,
			&e.StockNum, &e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose,
		); err != nil {
			return nil, err
		}
		entities = append(entities, &e)
	}
	return entities, nil
}