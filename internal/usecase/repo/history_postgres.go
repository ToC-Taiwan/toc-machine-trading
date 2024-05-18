package repo

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/pkg/postgres"
)

// history -.
type history struct {
	*postgres.Postgres
}

// NewHistory -.
func NewHistory(pg *postgres.Postgres) HistoryRepo {
	return &history{pg}
}

// InsertHistoryCloseArr -.
func (r *history) InsertHistoryCloseArr(ctx context.Context, t []*entity.StockHistoryClose) error {
	split := [][]*entity.StockHistoryClose{}
	if len(t) > batchSize {
		count := len(t)/batchSize + 1
		for i := 0; i < count; i++ {
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(t) {
				end = len(t)
			}
			if start != end {
				split = append(split, t[start:end])
			}
		}
	} else {
		split = append(split, t)
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	for _, s := range split {
		builder := r.Builder.Insert(tableNameHistoryStockClose).Columns("date, stock_num, close")
		for _, d := range s {
			builder = builder.Values(d.Date, d.StockNum, d.Close)
		}
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}

	return nil
}

// DeleteHistoryCloseByStockAndDate -.
func (r *history) DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.Delete(tableNameHistoryStockClose).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where(squirrel.Eq{"date": date})
	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

// QueryMutltiStockCloseByDate -.
func (r *history) QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.StockHistoryClose, error) {
	sql, args, err := r.Builder.
		Select("date, stock_num, close, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameHistoryStockClose).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where(squirrel.Eq{"date": date}).
		Join("basic_stock ON history_stock_close.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	closeMap := make(map[string]*entity.StockHistoryClose)
	for rows.Next() {
		e := entity.StockHistoryClose{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.Date, &e.StockNum, &e.Close,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate,
		); err != nil {
			return nil, err
		}
		closeMap[e.StockNum] = &e
	}
	return closeMap, nil
}

// InsertHistoryTickArr -.
func (r *history) InsertHistoryTickArr(ctx context.Context, t []*entity.StockHistoryTick) error {
	var split [][]*entity.StockHistoryTick
	if len(t) > batchSize {
		count := len(t)/batchSize + 1
		for i := 0; i < count; i++ {
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(t) {
				end = len(t)
			}
			if start != end {
				split = append(split, t[start:end])
			}
		}
	} else {
		split = append(split, t)
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	for _, s := range split {
		builder := r.Builder.Insert(tableNameHistoryStockTick).Columns("stock_num, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume")
		for _, v := range s {
			builder = builder.Values(v.StockNum, v.TickTime, v.Close, v.TickType, v.Volume, v.BidPrice, v.BidVolume, v.AskPrice, v.AskVolume)
		}

		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// DeleteHistoryTickByStockAndDate -.
func (r *history) DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.Delete(tableNameHistoryStockTick).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where(squirrel.GtOrEq{"tick_time": date}).
		Where(squirrel.Lt{"tick_time": date.AddDate(0, 0, 1)})
	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

// QueryMultiStockTickArrByDate -.
func (r *history) QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryTick, error) {
	sql, args, err := r.Builder.
		Select("stock_num, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameHistoryStockTick).
		Where(squirrel.GtOrEq{"tick_time": date}).
		Where(squirrel.Lt{"tick_time": date.AddDate(0, 0, 1)}).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		OrderBy("tick_time ASC").
		Join("basic_stock ON history_stock_tick.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*entity.StockHistoryTick)
	for rows.Next() {
		e := entity.StockHistoryTick{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.StockNum, &e.TickTime, &e.Close, &e.TickType, &e.Volume, &e.BidPrice, &e.BidVolume, &e.AskPrice, &e.AskVolume,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate,
		); err != nil {
			return nil, err
		}
		result[e.StockNum] = append(result[e.StockNum], &e)
	}
	return result, nil
}

// InsertHistoryKbarArr -.
func (r *history) InsertHistoryKbarArr(ctx context.Context, t []*entity.StockHistoryKbar) error {
	var split [][]*entity.StockHistoryKbar
	if len(t) > batchSize {
		count := len(t)/batchSize + 1
		for i := 0; i < count; i++ {
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(t) {
				end = len(t)
			}
			if start != end {
				split = append(split, t[start:end])
			}
		}
	} else {
		split = append(split, t)
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	for _, s := range split {
		builder := r.Builder.Insert(tableNameHistoryStockKbar).Columns("stock_num, kbar_time, open, high, low, close, volume")
		for _, v := range s {
			builder = builder.Values(v.StockNum, v.KbarTime, v.Open, v.High, v.Low, v.Close, v.Volume)
		}

		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// DeleteHistoryKbarByStockAndDate -.
func (r *history) DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.Delete(tableNameHistoryStockKbar).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where(squirrel.GtOrEq{"kbar_time": date}).
		Where(squirrel.Lt{"kbar_time": date.AddDate(0, 0, 1)})
	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

// QueryMultiStockKbarArrByDate -.
func (r *history) QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryKbar, error) {
	sql, args, err := r.Builder.
		Select("stock_num, kbar_time, open, high, low, close, volume, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameHistoryStockKbar).
		Where(squirrel.GtOrEq{"kbar_time": date}).
		Where(squirrel.Lt{"kbar_time": date.AddDate(0, 0, 1)}).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		OrderBy("kbar_time ASC").
		Join("basic_stock ON history_stock_kbar.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*entity.StockHistoryKbar)
	for rows.Next() {
		e := entity.StockHistoryKbar{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.StockNum, &e.KbarTime, &e.Open, &e.High, &e.Low, &e.Close, &e.Volume,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate,
		); err != nil {
			return nil, err
		}
		result[e.StockNum] = append(result[e.StockNum], &e)
	}
	return result, nil
}

// InsertQuaterMA -.
func (r *history) InsertQuaterMA(ctx context.Context, t *entity.StockHistoryAnalyze) error {
	dbQuaterMA, err := r.QueryAllQuaterMAByStockNum(ctx, t.StockNum)
	if err != nil {
		return err
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	if _, ok := dbQuaterMA[t.Date]; !ok {
		builder := r.Builder.Insert(tableNameHistoryStockAnalyze).Columns("date, stock_num, quater_ma").Values(t.Date, t.StockNum, t.QuaterMA)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryAllQuaterMAByStockNum -.
func (r *history) QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.StockHistoryAnalyze, error) {
	sql, args, err := r.Builder.
		Select("date, stock_num, quater_ma, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameHistoryStockAnalyze).
		Where(squirrel.Eq{"stock_num": stockNum}).
		OrderBy("date ASC").
		Join("basic_stock ON history_stock_analyze.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[time.Time]*entity.StockHistoryAnalyze)
	for rows.Next() {
		e := entity.StockHistoryAnalyze{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.Date, &e.StockNum, &e.QuaterMA, &e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result[e.Date] = &e
	}
	return result, nil
}
