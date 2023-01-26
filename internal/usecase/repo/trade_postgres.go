// Package repo package repo
package repo

import (
	"context"
	"errors"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4"
)

type TradeRepo struct {
	*postgres.Postgres
}

func NewTrade(pg *postgres.Postgres) usecase.TradeRepo {
	return &TradeRepo{pg}
}

// InsertOrUpdateOrderByOrderID -.
func (r *TradeRepo) InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error {
	dbOrder, err := r.QueryStockOrderByID(ctx, t.OrderID)
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
		builder := r.Builder.Insert(tableNameTradeStockOrder).Columns("manual, group_id, order_id, status, order_time, tick_time, stock_num, action, price, quantity, trade_time")
		builder = builder.Values(t.Manual, t.GroupID, t.OrderID, t.Status, t.OrderTime, t.TickTime, t.StockNum, t.Action, t.Price, t.Quantity, t.TradeTime)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else if !cmp.Equal(t, dbOrder) {
		builder := r.Builder.
			Update(tableNameTradeStockOrder).
			Set("manual", t.Manual).
			Set("group_id", t.GroupID).
			Set("order_id", t.OrderID).
			Set("status", t.Status).
			Set("order_time", t.OrderTime).
			Set("tick_time", t.TickTime).
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

// QueryStockOrderByID -.
func (r *TradeRepo) QueryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error) {
	sql, arg, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeStockOrder).
		Where(squirrel.Eq{"order_id": orderID}).
		Join("basic_stock ON trade_stock_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.StockOrder{Stock: new(entity.Stock)}
	if err := row.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
		&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllStockOrderByDate -.
func (r *TradeRepo) QueryAllStockOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.StockOrder, error) {
	sql, arg, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeStockOrder).
		Where(squirrel.GtOrEq{"order_time": timeRange[0]}).
		Where(squirrel.Lt{"order_time": timeRange[1]}).
		OrderBy("order_time ASC").
		Join("basic_stock ON trade_stock_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.StockOrder
	for rows.Next() {
		e := entity.StockOrder{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryAllStockOrder -.
func (r *TradeRepo) QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	sql, _, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, stock_num, action, price, quantity, trade_time, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeStockOrder).
		Join("basic_stock ON trade_stock_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.StockOrder
	for rows.Next() {
		e := entity.StockOrder{Stock: new(entity.Stock)}
		if err := rows.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.StockNum, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// InsertOrUpdateStockTradeBalance -.
func (r *TradeRepo) InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error {
	dbTradeBalance, err := r.QueryStockTradeBalanceByDate(ctx, t.TradeDay)
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
		builder := r.Builder.Insert(tableNameTradeStockBalance).Columns("trade_count, forward, reverse, original_balance, discount, total, trade_day")
		builder = builder.Values(t.TradeCount, t.Forward, t.Reverse, t.OriginalBalance, t.Discount, t.Total, t.TradeDay)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameTradeStockBalance).
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

// QueryStockTradeBalanceByDate -.
func (r *TradeRepo) QueryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.StockTradeBalance, error) {
	sql, arg, err := r.Builder.
		Select("trade_count, forward, reverse, original_balance, discount, total, trade_day").
		From(tableNameTradeStockBalance).
		Where(squirrel.Eq{"trade_day": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.StockTradeBalance{}
	if err := row.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.OriginalBalance, &e.Discount, &e.Total, &e.TradeDay); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllStockTradeBalance -.
func (r *TradeRepo) QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error) {
	sql, _, err := r.Builder.
		Select("trade_count, forward, reverse, original_balance, discount, total, trade_day").
		From(tableNameTradeStockBalance).OrderBy("trade_day ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.StockTradeBalance
	for rows.Next() {
		e := entity.StockTradeBalance{}
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

func (r *TradeRepo) QueryLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error) {
	sql, arg, err := r.Builder.
		Select("trade_count, forward, reverse, original_balance, discount, total, trade_day").
		From(tableNameTradeStockBalance).
		OrderBy("trade_day DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.StockTradeBalance{}
	if err := row.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.OriginalBalance, &e.Discount, &e.Total, &e.TradeDay); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllFutureTradeBalance -.
func (r *TradeRepo) QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error) {
	sql, _, err := r.Builder.
		Select("trade_count, forward, reverse, total, trade_day").
		From(tableNameFutureTradeBalance).OrderBy("trade_day ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.FutureTradeBalance
	for rows.Next() {
		e := entity.FutureTradeBalance{}
		if err := rows.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.Total, &e.TradeDay); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *TradeRepo) QueryLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error) {
	sql, arg, err := r.Builder.
		Select("trade_count, forward, reverse, total, trade_day").
		From(tableNameFutureTradeBalance).
		OrderBy("trade_day DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.FutureTradeBalance{}
	if err := row.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.Total, &e.TradeDay); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryFutureOrderByID -.
func (r *TradeRepo) QueryFutureOrderByID(ctx context.Context, orderID string) (*entity.FutureOrder, error) {
	sql, arg, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, trade_future_order.code, action, price, quantity, trade_time, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameTradeFutureOrder).
		Where(squirrel.Eq{"order_id": orderID}).
		Join("basic_future ON trade_future_order.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.FutureOrder{Future: new(entity.Future)}
	if err := row.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.Code, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
		&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// InsertOrUpdateFutureOrderByOrderID -.
func (r *TradeRepo) InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error {
	dbOrder, err := r.QueryFutureOrderByID(ctx, t.OrderID)
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
		builder := r.Builder.Insert(tableNameTradeFutureOrder).Columns("manual, group_id, order_id, status, order_time, tick_time, code, action, price, quantity, trade_time")
		builder = builder.Values(t.Manual, t.GroupID, t.OrderID, t.Status, t.OrderTime, t.TickTime, t.Code, t.Action, t.Price, t.Quantity, t.TradeTime)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else if !cmp.Equal(t, dbOrder) {
		builder := r.Builder.
			Update(tableNameTradeFutureOrder).
			Set("manual", t.Manual).
			Set("group_id", t.GroupID).
			Set("order_id", t.OrderID).
			Set("status", t.Status).
			Set("order_time", t.OrderTime).
			Set("tick_time", t.TickTime).
			Set("code", t.Code).
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

func (r *TradeRepo) QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
	sql, _, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, trade_future_order.code, action, price, quantity, trade_time, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameTradeFutureOrder).
		OrderBy("order_time ASC").
		Join("basic_future ON trade_future_order.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.FutureOrder
	for rows.Next() {
		e := entity.FutureOrder{Future: new(entity.Future)}
		if err := rows.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.Code, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryAllFutureOrderByDate -.
func (r *TradeRepo) QueryAllFutureOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.FutureOrder, error) {
	sql, arg, err := r.Builder.
		Select("manual, group_id, order_id, status, order_time, tick_time, trade_future_order.code, action, price, quantity, trade_time, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameTradeFutureOrder).
		Where(squirrel.GtOrEq{"order_time": timeRange[0]}).
		Where(squirrel.Lt{"order_time": timeRange[1]}).
		OrderBy("order_time ASC").
		Join("basic_future ON trade_future_order.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.FutureOrder
	for rows.Next() {
		e := entity.FutureOrder{Future: new(entity.Future)}
		if err := rows.Scan(&e.Manual, &e.GroupID, &e.OrderID, &e.Status, &e.OrderTime, &e.TickTime, &e.Code, &e.Action, &e.Price, &e.Quantity, &e.TradeTime,
			&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryFutureTradeBalanceByDate -.
func (r *TradeRepo) QueryFutureTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.FutureTradeBalance, error) {
	sql, arg, err := r.Builder.
		Select("trade_count, forward, reverse, total, trade_day").
		From(tableNameFutureTradeBalance).
		Where(squirrel.Eq{"trade_day": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.FutureTradeBalance{}
	if err := row.Scan(&e.TradeCount, &e.Forward, &e.Reverse, &e.Total, &e.TradeDay); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// InsertOrUpdateFutureTradeBalance -.
func (r *TradeRepo) InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.FutureTradeBalance) error {
	dbTradeBalance, err := r.QueryFutureTradeBalanceByDate(ctx, t.TradeDay)
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
		builder := r.Builder.Insert(tableNameFutureTradeBalance).Columns("trade_count, forward, reverse, total, trade_day")
		builder = builder.Values(t.TradeCount, t.Forward, t.Reverse, t.Total, t.TradeDay)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameFutureTradeBalance).
			Set("trade_count", t.TradeCount).
			Set("forward", t.Forward).
			Set("reverse", t.Reverse).
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
