package repo

import (
	"context"
	"errors"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
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
	log.Infof("InsertHistoryCloseArr -> Count: %d", len(t))
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
		Where(squirrel.Eq{"date": date}).
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

// InsertHistoryTickArr -.
func (r *HistoryRepo) InsertHistoryTickArr(ctx context.Context, t []*entity.HistoryTick) error {
	log.Infof("InsertHistoryTickArr -> Count: %d", len(t))
	var split [][]*entity.HistoryTick
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
		builder := r.Builder.Insert(tableNameHistoryTick).Columns("stock_num, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume")
		for _, v := range s {
			builder = builder.Values(v.StockNum, v.TickTime, v.Close, v.TickType, v.Volume, v.BidPrice, v.BidVolume, v.AskPrice, v.AskVolume)
		}

		if sql, args, err := builder.ToSql(); err != nil {
			return err
		} else if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// CheckHistoryTickExist -.
func (r *HistoryRepo) CheckHistoryTickExist(ctx context.Context, stockNum string, date time.Time) (bool, error) {
	sql, args, err := r.Builder.
		Select("stock_num").
		From(tableNameHistoryTick).
		Where(squirrel.GtOrEq{"tick_time": date}).
		Where(squirrel.Lt{"tick_time": date.Add(time.Hour * 24)}).
		Where(squirrel.Eq{"stock_num": stockNum}).Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	row := r.Pool.QueryRow(ctx, sql, args...)
	var e string
	if err := row.Scan(&e); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	if e != "" {
		return true, nil
	}
	return false, nil
}

// InsertHistoryKbarArr -.
func (r *HistoryRepo) InsertHistoryKbarArr(ctx context.Context, t []*entity.HistoryKbar) error {
	log.Infof("InsertHistoryKbarArr -> Count: %d", len(t))
	var split [][]*entity.HistoryKbar
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
		builder := r.Builder.Insert(tableNameHistoryKbar).Columns("stock_num, kbar_time, open, high, low, close, volume")
		for _, v := range s {
			builder = builder.Values(v.StockNum, v.KbarTime, v.Open, v.High, v.Low, v.Close, v.Volume)
		}

		if sql, args, err := builder.ToSql(); err != nil {
			return err
		} else if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// CheckHistoryKbarExist -.
func (r *HistoryRepo) CheckHistoryKbarExist(ctx context.Context, stockNum string, date time.Time) (bool, error) {
	sql, args, err := r.Builder.
		Select("stock_num").
		From(tableNameHistoryKbar).
		Where(squirrel.GtOrEq{"kbar_time": date}).
		Where(squirrel.Lt{"kbar_time": date.Add(time.Hour * 24)}).
		Where(squirrel.Eq{"stock_num": stockNum}).Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	row := r.Pool.QueryRow(ctx, sql, args...)
	var e string
	if err := row.Scan(&e); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	if e != "" {
		return true, nil
	}
	return false, nil
}
