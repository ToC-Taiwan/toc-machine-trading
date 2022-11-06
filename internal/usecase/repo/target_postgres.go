package repo

import (
	"context"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/postgres"
)

// TargetRepo -.
type TargetRepo struct {
	*postgres.Postgres
}

// NewTarget -.
func NewTarget(pg *postgres.Postgres) *TargetRepo {
	return &TargetRepo{pg}
}

// InsertOrUpdateTargetArr -.
func (r *TargetRepo) InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error {
	inDBTargets, err := r.QueryTargetsByTradeDay(ctx, t[0].TradeDay)
	if err != nil {
		return err
	}

	inDBTargetsMap := make(map[string]*entity.StockTarget)
	for _, v := range inDBTargets {
		inDBTargetsMap[v.StockNum] = v
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	var insert int
	builder := r.Builder.Insert(tableNameTarget).Columns("stock_num, trade_day, rank, volume")
	for _, v := range t {
		if _, ok := inDBTargetsMap[v.StockNum]; !ok {
			insert++
			builder = builder.Values(v.StockNum, v.TradeDay, v.Rank, v.Volume)
		} else {
			b := r.Builder.
				Update(tableNameTarget).
				Set("stock_num", v.StockNum).
				Set("trade_day", v.TradeDay).
				Set("rank", v.Rank).
				Set("volume", v.Volume).
				Where("stock_num = ?", v.StockNum).
				Where("trade_day = ?", v.TradeDay)
			if sql, args, err = b.ToSql(); err != nil {
				return err
			} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
				return err
			}
		}
	}

	if insert != 0 {
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryTargetsByTradeDay -.
func (r *TargetRepo) QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error) {
	sql, args, err := r.Builder.
		Select("id, rank, volume, trade_day, stock_num, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTarget).
		Where("trade_day = ?", tradeDay).
		OrderBy("rank ASC").
		Join("basic_stock ON basic_targets.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.StockTarget, 0, 256)
	for rows.Next() {
		e := entity.StockTarget{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.ID, &e.Rank, &e.Volume, &e.TradeDay,
			&e.StockNum, &e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate,
		); err != nil {
			return nil, err
		}
		entities = append(entities, &e)
	}
	return entities, nil
}
