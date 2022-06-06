// Package repo package repo
package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/logger"
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
	inDBStockMap := make(map[string]*entity.Stock)
	for _, s := range inDBStock {
		inDBStockMap[s.Number] = s
	}

	var insert, update int
	builder := r.Builder.Insert("basic_stock").Columns("number, name, exchange, category, day_trade, last_close")
	for _, v := range t {
		if _, ok := inDBStockMap[v.Number]; !ok {
			insert++
			builder = builder.Values(v.Number, v.Name, v.Exchange, v.Category, v.DayTrade, v.LastClose)
		} else if !cmp.Equal(v, inDBStockMap[v.Number]) {
			update++
			builder := r.Builder.
				Update("basic_stock").
				Set("number", v.Number).
				Set("name", v.Name).
				Set("exchange", v.Exchange).
				Set("category", v.Category).
				Set("day_trade", v.DayTrade).
				Set("last_close", v.LastClose).
				Where("number = ?", v.Number)
			sql, args, updateErr := builder.ToSql()
			if updateErr != nil {
				return err
			}
			_, err = r.Pool.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
	}

	if insert != 0 {
		sql, args, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	logger.Get().WithFields(map[string]interface{}{
		"Update": update,
		"Exist":  len(t) - update - insert,
		"Insert": insert,
	}).Info("InserOrUpdatetStockArr")
	return nil
}

// QueryAllStock -.
func (r *BasicRepo) QueryAllStock(ctx context.Context) ([]*entity.Stock, error) {
	sql, _, err := r.Builder.
		Select("*").
		From("basic_stock").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.Stock, 0, 2048)
	for rows.Next() {
		e := &entity.Stock{}
		err = rows.Scan(
			&e.Number,
			&e.Name,
			&e.Exchange,
			&e.Category,
			&e.DayTrade,
			&e.LastClose,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}

// InserOrUpdatetCalendarDateArr -.
func (r *BasicRepo) InserOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error {
	inDBCalendar, err := r.QueryAllTradeDay(ctx)
	if err != nil {
		return err
	}
	inDBCalendarMap := make(map[string]*entity.CalendarDate)
	for _, s := range inDBCalendar {
		inDBCalendarMap[s.Date.String()] = s
	}

	var insert, update int
	builder := r.Builder.Insert("basic_calendar").Columns("date, is_trade_day")
	for _, v := range t {
		if _, ok := inDBCalendarMap[v.Date.String()]; !ok {
			insert++
			builder = builder.Values(v.Date, v.IsTradeDay)
		} else if !cmp.Equal(v, inDBCalendarMap[v.Date.String()]) {
			update++
			builder := r.Builder.
				Update("basic_calendar").
				Set("date", v.Date).
				Set("is_trade_day", v.IsTradeDay).
				Where("date = ?", v.Date)
			sql, args, updateErr := builder.ToSql()
			if updateErr != nil {
				return err
			}
			_, err = r.Pool.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
	}

	if insert != 0 {
		sql, args, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	logger.Get().WithFields(map[string]interface{}{
		"Update": update,
		"Exist":  len(t) - update - insert,
		"Insert": insert,
	}).Info("InserOrUpdatetCalendarDateArr")
	return nil
}

// QueryAllTradeDay -.
func (r *BasicRepo) QueryAllTradeDay(ctx context.Context) ([]*entity.CalendarDate, error) {
	sql, _, err := r.Builder.
		Select("*").
		From("basic_calendar").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*entity.CalendarDate, 0, 1024)
	for rows.Next() {
		e := &entity.CalendarDate{}
		err = rows.Scan(
			&e.Date,
			&e.IsTradeDay,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}
