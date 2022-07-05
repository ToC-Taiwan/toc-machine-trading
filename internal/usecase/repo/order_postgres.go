// Package repo package repo
package repo

import (
	"context"
	"errors"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4"
)

// OrderRepo -.
type OrderRepo struct {
	*postgres.Postgres
}

// NewOrder -.
func NewOrder(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

// InsertOrUpdateOrder -.
func (r *OrderRepo) InsertOrUpdateOrder(ctx context.Context, t *entity.Order) error {
	dbOrder, err := r.QueryOrderByID(ctx, t.OrderID)
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

	if dbOrder == nil {
		builder := r.Builder.Insert(tableNameTradeOrder).Columns("uuid ,order_id, status, order_time, stock_num, action, price, quantity, trade_time")
		builder = builder.Values(t.UUID, t.OrderID, t.Status, t.OrderTime, t.StockNum, t.Action, t.Price, t.Quantity, t.TradeTime)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else if !cmp.Equal(t, dbOrder) {
		builder := r.Builder.
			Update(tableNameTradeOrder).
			Set("uuid", t.UUID).
			Set("order_id", t.OrderID).
			Set("status", t.Status).
			Set("order_time", t.OrderTime).
			Set("stock_num", t.StockNum).
			Set("action", t.Action).
			Set("price", t.Price).
			Set("quantity", t.Quantity).
			Set("trade_time", t.TradeTime).
			Where(squirrel.Eq{"order_id": t.OrderID})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryOrderByID -.
func (r *OrderRepo) QueryOrderByID(ctx context.Context, orderID string) (*entity.Order, error) {
	sql, arg, err := r.Builder.
		Select("uuid, order_id, status, order_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeOrder).
		Where(squirrel.Eq{"order_id": orderID}).
		Join("basic_stock ON trade_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.Order{Stock: new(entity.Stock)}
	if err := row.Scan(&e.UUID, &e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
		&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllOrderByDate -.
func (r *OrderRepo) QueryAllOrderByDate(ctx context.Context, date time.Time) ([]*entity.Order, error) {
	sql, arg, err := r.Builder.
		Select("uuid, order_id, status, order_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeOrder).
		Where(squirrel.GtOrEq{"order_time": date}).
		Where(squirrel.Lt{"order_time": date.AddDate(0, 0, 1)}).
		Join("basic_stock ON trade_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Order
	for rows.Next() {
		e := entity.Order{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.UUID, &e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryAllOrder -.
func (r *OrderRepo) QueryAllOrder(ctx context.Context) ([]*entity.Order, error) {
	sql, _, err := r.Builder.
		Select("uuid ,order_id, status, order_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeOrder).
		Join("basic_stock ON trade_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Order
	for rows.Next() {
		e := entity.Order{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.UUID, &e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// InsertOrUpdateTradeBalance -.
func (r *OrderRepo) InsertOrUpdateTradeBalance(ctx context.Context, t *entity.TradeBalance) error {
	dbTradeBalance, err := r.QueryTradeBalanceByDate(ctx, t.TradeDay)
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

	if dbTradeBalance == nil {
		builder := r.Builder.Insert(tableNameTradeBalance).Columns("trade_count, forward, reverse, original_balance, discount, total, trade_day")
		builder = builder.Values(t.TradeCount, t.Forward, t.Reverse, t.OriginalBalance, t.Discount, t.Total, t.TradeDay)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameTradeBalance).
			Set("trade_count", t.TradeCount).
			Set("forward", t.Forward).
			Set("reverse", t.Reverse).
			Set("original_balance", t.OriginalBalance).
			Set("discount", t.Discount).
			Set("total", t.Total).
			Set("trade_day", t.TradeDay).
			Where(squirrel.Eq{"trade_day": t.TradeDay})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// QueryTradeBalanceByDate -.
func (r *OrderRepo) QueryTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.TradeBalance, error) {
	sql, arg, err := r.Builder.
		Select("trade_count, forward, reverse, original_balance, discount, total, trade_day").
		From(tableNameTradeBalance).
		Where(squirrel.Eq{"trade_day": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.TradeBalance{}
	if err := row.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.OriginalBalance, &e.Discount, &e.Total, &e.TradeDay); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllTradeBalance -.
func (r *OrderRepo) QueryAllTradeBalance(ctx context.Context) ([]*entity.TradeBalance, error) {
	sql, _, err := r.Builder.
		Select("trade_count, forward, reverse, original_balance, discount, total, trade_day").
		From(tableNameTradeBalance).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.TradeBalance
	for rows.Next() {
		e := entity.TradeBalance{}
		if err := rows.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.OriginalBalance, &e.Discount, &e.Total, &e.TradeDay); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}
