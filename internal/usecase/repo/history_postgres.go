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
	log.Infof("InsertHistoryCloseArr -> Count: %d", len(t))
	split := [][]*entity.HistoryClose{}
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
		builder := r.Builder.Insert(tableNameHistoryClose).Columns("date, stock_num, close")
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

// QueryMutltiStockCloseByDate -.
func (r *HistoryRepo) QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.HistoryClose, error) {
	sql, args, err := r.Builder.
		Select("date, stock_num, close, number, name, exchange, category, day_trade, last_close").
		From(tableNameHistoryClose).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		Where(squirrel.Eq{"date": date}).
		Join("basic_stock ON history_close.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
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
		builder := r.Builder.Insert(tableNameHistoryTick).Columns("stock_num, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume")
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

// QueryMultiStockTickArrByDate -.
func (r *HistoryRepo) QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.HistoryTick, error) {
	sql, args, err := r.Builder.
		Select("stock_num, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume, number, name, exchange, category, day_trade, last_close").
		From(tableNameHistoryTick).
		Where(squirrel.GtOrEq{"tick_time": date}).
		Where(squirrel.Lt{"tick_time": date.AddDate(0, 0, 1)}).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		OrderBy("tick_time ASC").
		Join("basic_stock ON history_tick.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*entity.HistoryTick)
	for rows.Next() {
		e := entity.HistoryTick{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.StockNum, &e.TickTime, &e.Close, &e.TickType, &e.Volume, &e.BidPrice, &e.BidVolume, &e.AskPrice, &e.AskVolume,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose,
		); err != nil {
			return nil, err
		}
		result[e.StockNum] = append(result[e.StockNum], &e)
	}
	return result, nil
}

// InsertHistoryKbarArr -.
func (r *HistoryRepo) InsertHistoryKbarArr(ctx context.Context, t []*entity.HistoryKbar) error {
	log.Infof("InsertHistoryKbarArr -> Count: %d", len(t))
	var split [][]*entity.HistoryKbar
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
		builder := r.Builder.Insert(tableNameHistoryKbar).Columns("stock_num, kbar_time, open, high, low, close, volume")
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

// QueryMultiStockKbarArrByDate -.
func (r *HistoryRepo) QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.HistoryKbar, error) {
	sql, args, err := r.Builder.
		Select("stock_num, kbar_time, open, high, low, close, volume, number, name, exchange, category, day_trade, last_close").
		From(tableNameHistoryKbar).
		Where(squirrel.GtOrEq{"kbar_time": date}).
		Where(squirrel.Lt{"kbar_time": date.AddDate(0, 0, 1)}).
		Where(squirrel.Eq{"stock_num": stockNumArr}).
		OrderBy("kbar_time ASC").
		Join("basic_stock ON history_kbar.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*entity.HistoryKbar)
	for rows.Next() {
		e := entity.HistoryKbar{Stock: new(entity.Stock)}
		if err := rows.Scan(
			&e.StockNum, &e.KbarTime, &e.Open, &e.High, &e.Low, &e.Close, &e.Volume,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose,
		); err != nil {
			return nil, err
		}
		result[e.StockNum] = append(result[e.StockNum], &e)
	}
	return result, nil
}

// InsertQuaterMA -.
func (r *HistoryRepo) InsertQuaterMA(ctx context.Context, t *entity.HistoryAnalyze) error {
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
		builder := r.Builder.Insert(tableNameHistoryAnalyze).Columns("date, stock_num, quater_ma").Values(t.Date, t.StockNum, t.QuaterMA)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryAllQuaterMAByStockNum -.
func (r *HistoryRepo) QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.HistoryAnalyze, error) {
	sql, args, err := r.Builder.
		Select("date, stock_num, quater_ma, number, name, exchange, category, day_trade, last_close").
		From(tableNameHistoryAnalyze).
		Where(squirrel.Eq{"stock_num": stockNum}).
		OrderBy("date ASC").
		Join("basic_stock ON history_analyze.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[time.Time]*entity.HistoryAnalyze)
	for rows.Next() {
		e := entity.HistoryAnalyze{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.Date, &e.StockNum, &e.QuaterMA, &e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose); err != nil {
			return nil, err
		}
		result[e.Date] = &e
	}
	return result, nil
}
