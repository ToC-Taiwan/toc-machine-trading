// Package repo package repo
package repo

import (
	"context"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"

	"github.com/google/go-cmp/cmp"
)

// BasicRepo -.
type BasicRepo struct {
	*postgres.Postgres
}

// NewBasic -.
func NewBasic(pg *postgres.Postgres) *BasicRepo {
	return &BasicRepo{pg}
}

// InserOrUpdatetStockArr -.
func (r *BasicRepo) InserOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error {
	inDBStock, err := r.QueryAllStock(ctx)
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
	builder := r.Builder.Insert(tableNameStock).Columns("number, name, exchange, category, day_trade, last_close")
	for _, v := range t {
		if _, ok := inDBStock[v.Number]; !ok {
			insert++
			builder = builder.Values(v.Number, v.Name, v.Exchange, v.Category, v.DayTrade, v.LastClose)
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
				Where("number = ?", v.Number)
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
	log.Infof("InserOrUpdatetStockArr -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

// QueryAllStock -.
func (r *BasicRepo) QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error) {
	sql, _, err := r.Builder.
		Select("number, name, exchange, category, day_trade, last_close").
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
		if err = rows.Scan(&e.Number, &e.Name, &e.Exchange, &e.Category, &e.DayTrade, &e.LastClose); err != nil {
			return nil, err
		}
		entities[e.Number] = &e
	}
	return entities, nil
}

// InserOrUpdatetCalendarDateArr -.
func (r *BasicRepo) InserOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error {
	inDBCalendar, err := r.QueryAllCalendar(ctx)
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
	log.Infof("InserOrUpdatetCalendarDateArr -> Exist: %d, Insert: %d, Update: %d", len(t)-update-insert, insert, update)
	return nil
}

// QueryAllCalendar -.
func (r *BasicRepo) QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error) {
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
