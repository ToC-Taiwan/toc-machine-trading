package repo

import (
	"context"
	"errors"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/postgres"

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
func (r *HistoryRepo) InsertHistoryCloseArr(ctx context.Context, t []*entity.StockHistoryClose) error {
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
func (r *HistoryRepo) DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
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
func (r *HistoryRepo) QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.StockHistoryClose, error) {
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
func (r *HistoryRepo) InsertHistoryTickArr(ctx context.Context, t []*entity.StockHistoryTick) error {
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
func (r *HistoryRepo) DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
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
func (r *HistoryRepo) QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryTick, error) {
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
func (r *HistoryRepo) InsertHistoryKbarArr(ctx context.Context, t []*entity.StockHistoryKbar) error {
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
func (r *HistoryRepo) DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
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
func (r *HistoryRepo) QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryKbar, error) {
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
func (r *HistoryRepo) InsertQuaterMA(ctx context.Context, t *entity.StockHistoryAnalyze) error {
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
func (r *HistoryRepo) QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.StockHistoryAnalyze, error) {
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

// InsertFutureHistoryTickArr -.
func (r *HistoryRepo) InsertFutureHistoryTickArr(ctx context.Context, t []*entity.FutureHistoryTick) error {
	var split [][]*entity.FutureHistoryTick
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
		builder := r.Builder.Insert(tableNameHistoryFutureTick).Columns("code, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume")
		for _, v := range s {
			builder = builder.Values(v.Code, v.TickTime, v.Close, v.TickType, v.Volume, v.BidPrice, v.BidVolume, v.AskPrice, v.AskVolume)
		}

		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryFutureHistoryTickArrByTime -.
func (r *HistoryRepo) QueryFutureHistoryTickArrByTime(ctx context.Context, code string, startTime, endTime time.Time) ([]*entity.FutureHistoryTick, error) {
	sql, args, err := r.Builder.
		Select("history_future_tick.code, tick_time, close, tick_type, volume, bid_price, bid_volume, ask_price, ask_volume, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameHistoryFutureTick).
		Where(squirrel.GtOrEq{"tick_time": startTime}).
		Where(squirrel.Lt{"tick_time": endTime}).
		Where(squirrel.Eq{"history_future_tick.code": code}).
		OrderBy("tick_time ASC").
		Join("basic_future ON history_future_tick.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*entity.FutureHistoryTick{}
	for rows.Next() {
		e := entity.FutureHistoryTick{Future: new(entity.Future)}
		if err := rows.Scan(
			&e.Code, &e.TickTime, &e.Close, &e.TickType, &e.Volume, &e.BidPrice, &e.BidVolume, &e.AskPrice, &e.AskVolume,
			&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate,
		); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *HistoryRepo) InsertFutureHistoryClose(ctx context.Context, c *entity.FutureHistoryClose) error {
	dbClose, err := r.QueryFutureHistoryCloseByDate(ctx, c.Code, c.Date)
	if err != nil {
		return err
	}

	if dbClose.Close != 0 {
		return nil
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.Insert(tableNameHistoryFutureClose).Columns("date, code, close").Values(c.Date, c.Code, c.Close)
	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *HistoryRepo) QueryFutureHistoryCloseByDate(ctx context.Context, code string, tradeDay time.Time) (*entity.FutureHistoryClose, error) {
	sql, args, err := r.Builder.
		Select("history_future_close.code, date, close, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameHistoryFutureClose).
		Where(squirrel.Eq{"history_future_close.code": code}).
		Where(squirrel.Eq{"date": tradeDay}).
		Join("basic_future ON history_future_close.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, args...)
	e := entity.FutureHistoryClose{Future: new(entity.Future)}
	if err := row.Scan(
		&e.Code, &e.Date, &e.Close,
		&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.FutureHistoryClose{}, nil
		}
		return nil, err
	}
	return &e, nil
}
