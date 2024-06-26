// Package repo package repo
package repo

import (
	"context"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/toc-taiwan/postgres"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
)

// basic -.
type basic struct {
	*postgres.Postgres
}

// NewBasic -.
func NewBasic(pg *postgres.Postgres) BasicRepo {
	return &basic{pg}
}

// InsertOrUpdatetStockArr -.
func (r *basic) InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error {
	inDBStock, err := r.queryAllStock(ctx)
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

	var insert, update int
	builder := r.Builder.Insert(tableNameStock).Columns("number, name, exchange, category, day_trade, last_close, update_date")
	for _, v := range t {
		if _, ok := inDBStock[v.Number]; !ok {
			insert++
			builder = builder.Values(v.Number, v.Name, v.Exchange, v.Category, v.DayTrade, v.LastClose, v.UpdateDate)
		} else if !cmp.Equal(v, inDBStock[v.Number]) {
			update++
			b := r.Builder.
				Update(tableNameStock).
				Set("number", v.Number).
				Set("name", v.Name).
				Set("exchange", v.Exchange).
				Set("category", v.Category).
				Set("day_trade", v.DayTrade).
				Set("last_close", v.LastClose).
				Set("update_date", v.UpdateDate).
				Where("number = ?", v.Number)
			sql, args, err = b.ToSql()
			if err != nil {
				return err
			} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
				return err
			}
		}
	}

	if insert != 0 {
		sql, args, err = builder.ToSql()
		if err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	r.Infof("Insert Stock -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

// UpdateAllStockDayTradeToNo -.
func (r *basic) UpdateAllStockDayTradeToNo(ctx context.Context) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.
		Update(tableNameStock).
		Set("day_trade", false)

	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

// queryAllStock -.
func (r *basic) queryAllStock(ctx context.Context) (map[string]*entity.Stock, error) {
	sql, _, err := r.Builder.
		Select("number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameStock).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make(map[string]*entity.Stock)
	for rows.Next() {
		e := entity.Stock{}
		if err = rows.Scan(&e.Number, &e.Name, &e.Exchange, &e.Category, &e.DayTrade, &e.LastClose, &e.UpdateDate); err != nil {
			return nil, err
		}
		entities[e.Number] = &e
	}
	return entities, nil
}

// InsertOrUpdatetCalendarDateArr -.
func (r *basic) InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error {
	inDBCalendar, err := r.queryAllCalendar(ctx)
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

	var insert, update int
	builder := r.Builder.Insert(tableNameCalendar).Columns("date, is_trade_day")
	for _, v := range t {
		if _, ok := inDBCalendar[v.Date]; !ok {
			insert++
			builder = builder.Values(v.Date, v.IsTradeDay)
		} else if !cmp.Equal(v, inDBCalendar[v.Date]) {
			update++
			b := r.Builder.
				Update(tableNameCalendar).
				Set("date", v.Date).
				Set("is_trade_day", v.IsTradeDay).
				Where("date = ?", v.Date)
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
	r.Infof("Insert Calendar -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

// queryAllCalendar -.
func (r *basic) queryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error) {
	sql, _, err := r.Builder.
		Select("date, is_trade_day").
		From(tableNameCalendar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make(map[time.Time]*entity.CalendarDate)
	for rows.Next() {
		e := entity.CalendarDate{}
		if err = rows.Scan(&e.Date, &e.IsTradeDay); err != nil {
			return nil, err
		}
		entities[e.Date] = &e
	}
	return entities, nil
}

// InsertOrUpdatetFutureArr -.
func (r *basic) InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error {
	inDBFuture, err := r.queryAllFuture(ctx)
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

	var insert, update int
	builder := r.Builder.Insert(tableNameFuture).Columns("code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date")
	for _, v := range t {
		if _, ok := inDBFuture[v.Code]; !ok {
			insert++
			builder = builder.Values(v.Code, v.Symbol, v.Name, v.Category, v.DeliveryMonth, v.DeliveryDate, v.UnderlyingKind, v.Unit, v.LimitUp, v.LimitDown, v.Reference, v.UpdateDate)
		} else if !cmp.Equal(v, inDBFuture[v.Code]) {
			update++
			b := r.Builder.
				Update(tableNameFuture).
				Set("code", v.Code).
				Set("symbol", v.Symbol).
				Set("name", v.Name).
				Set("category", v.Category).
				Set("delivery_month", v.DeliveryMonth).
				Set("delivery_date", v.DeliveryDate).
				Set("underlying_kind", v.UnderlyingKind).
				Set("unit", v.Unit).
				Set("limit_up", v.LimitUp).
				Set("limit_down", v.LimitDown).
				Set("reference", v.Reference).
				Set("update_date", v.UpdateDate).
				Where("code = ?", v.Code)
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
	r.Infof("Insert Future -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

// queryAllFuture -.
func (r *basic) queryAllFuture(ctx context.Context) (map[string]*entity.Future, error) {
	sql, _, err := r.Builder.
		Select("code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		OrderBy("delivery_date ASC").
		From(tableNameFuture).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make(map[string]*entity.Future)
	for rows.Next() {
		e := entity.Future{}
		if err = rows.Scan(&e.Code, &e.Symbol, &e.Name, &e.Category, &e.DeliveryMonth, &e.DeliveryDate, &e.UnderlyingKind, &e.Unit, &e.LimitUp, &e.LimitDown, &e.Reference, &e.UpdateDate); err != nil {
			return nil, err
		}
		entities[e.Code] = &e
	}
	return entities, nil
}

func (r *basic) InsertOrUpdatetOptionArr(ctx context.Context, t []*entity.Option) error {
	inDBOption, err := r.queryAllOption(ctx)
	if err != nil {
		return err
	}
	var insert, update int
	var sql []string
	var args [][]interface{}
	for _, v := range t {
		if _, ok := inDBOption[v.Code]; !ok {
			insert++
			builder := r.Builder.Insert(tableNameOption).Columns("code, symbol, name, category, delivery_month, delivery_date, strike_price, option_right, underlying_kind, unit, limit_up, limit_down, reference, update_date")
			builder = builder.Values(v.Code, v.Symbol, v.Name, v.Category, v.DeliveryMonth, v.DeliveryDate, v.StrikePrice, v.OptionRight, v.UnderlyingKind, v.Unit, v.LimitUp, v.LimitDown, v.Reference, v.UpdateDate)
			if sqlCommand, argsCommand, ierr := builder.ToSql(); ierr != nil {
				return ierr
			} else if len(argsCommand) > 0 {
				sql = append(sql, sqlCommand)
				args = append(args, argsCommand)
			}
		} else if !cmp.Equal(v, inDBOption[v.Code]) {
			update++
			b := r.Builder.
				Update(tableNameOption).
				Set("code", v.Code).
				Set("symbol", v.Symbol).
				Set("name", v.Name).
				Set("category", v.Category).
				Set("delivery_month", v.DeliveryMonth).
				Set("delivery_date", v.DeliveryDate).
				Set("strike_price", v.StrikePrice).
				Set("option_right", v.OptionRight).
				Set("underlying_kind", v.UnderlyingKind).
				Set("unit", v.Unit).
				Set("limit_up", v.LimitUp).
				Set("limit_down", v.LimitDown).
				Set("reference", v.Reference).
				Set("update_date", v.UpdateDate).
				Where("code = ?", v.Code)
			if sqlCommand, argsCommand, e := b.ToSql(); e != nil {
				return e
			} else {
				sql = append(sql, sqlCommand)
				args = append(args, argsCommand)
			}
		}
	}

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	var sqlErr error
	defer r.EndTransaction(tx, sqlErr)
	for i, v := range sql {
		_, sqlErr = tx.Exec(ctx, v, args[i]...)
		if sqlErr != nil {
			return sqlErr
		}
	}
	r.Infof("Insert Option -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

func (r *basic) queryAllOption(ctx context.Context) (map[string]*entity.Option, error) {
	sql, _, err := r.Builder.
		Select("code, symbol, name, category, delivery_month, delivery_date, strike_price, option_right, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		OrderBy("delivery_date ASC").
		From(tableNameOption).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make(map[string]*entity.Option)
	for rows.Next() {
		e := entity.Option{}
		if err = rows.Scan(&e.Code, &e.Symbol, &e.Name, &e.Category, &e.DeliveryMonth, &e.DeliveryDate, &e.StrikePrice, &e.OptionRight, &e.UnderlyingKind, &e.Unit, &e.LimitUp, &e.LimitDown, &e.Reference, &e.UpdateDate); err != nil {
			return nil, err
		}
		entities[e.Code] = &e
	}
	return entities, nil
}
