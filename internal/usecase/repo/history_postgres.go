package repo

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

// HistoryRepo -.
type HistoryRepo struct {
	*postgres.Postgres
}

// NewHistory -.
func NewHistory(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg}
}

// InsertHistoryCloseArr -.
func (r *HistoryRepo) InsertHistoryCloseArr(ctx context.Context, t []*entity.HistoryClose) error {
	var split [][]*entity.HistoryClose
	count := len(t)/batchSize + 1
	for i := 0; i < count; i++ {
		if i == count-1 {
			if l := len(t[batchSize*i:]); l != 0 {
				split = append(split, t[batchSize*i:])
			}
		} else {
			split = append(split, t[batchSize*i:batchSize*(i+1)])
		}
	}

	for _, s := range split {
		builder := r.Builder.Insert(tableNameHistoryClose).Columns("date, stock_num, close")
		for _, v := range s {
			builder = builder.Values(v.Date, v.StockNum, v.Close)
		}

		if sql, args, err := builder.ToSql(); err != nil {
			return err
		} else if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryHistoryCloseByMutltiStockNumDate -.
func (r *HistoryRepo) QueryHistoryCloseByMutltiStockNumDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.HistoryClose, error) {
	sql, args, err := r.Builder.
		Select("date, stock_num, close, number, name, exchange, category, day_trade, last_close").
		From(tableNameHistoryClose).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where("date = ?", date).
		Join("basic_stock ON history_close.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	closeMap := make(map[string]*entity.HistoryClose)
	for rows.Next() {
		e := entity.HistoryClose{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.Date, &e.StockNum, &e.Close,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose,
		); err != nil {
			return nil, err
		}
		closeMap[e.StockNum] = &e
	}
	return closeMap, nil
}
