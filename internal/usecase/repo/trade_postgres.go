// Package repo package repo
package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4"
)

type trade struct {
	*postgres.Postgres
}

func NewTrade(pg *postgres.Postgres) TradeRepo {
	return &trade{pg}
}

// InsertOrUpdateOrderByOrderID -.
func (r *trade) InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error {
	dbOrder, err := r.queryStockOrderByID(ctx, t.OrderID)
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
		builder := r.Builder.
			Insert(tableNameTradeStockOrder).
			Columns("order_id, status, order_time, stock_num, action, price, lot, share")
		builder = builder.Values(t.OrderID, t.Status, t.OrderTime, t.StockNum, t.Action, t.Price, t.Lot, t.Share)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else if !cmp.Equal(t, dbOrder) {
		builder := r.Builder.
			Update(tableNameTradeStockOrder).
			Set("order_id", t.OrderID).
			Set("status", t.Status).
			Set("order_time", t.OrderTime).
			Set("stock_num", t.StockNum).
			Set("action", t.Action).
			Set("price", t.Price).
			Set("lot", t.Lot).
			Set("share", t.Share).
			Where(squirrel.Eq{"order_id": t.OrderID})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

// queryStockOrderByID -.
func (r *trade) queryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error) {
	sql, arg, err := r.Builder.
		Select("order_id, status, order_time, stock_num, action, price, lot, share, number, name, exchange, category, day_trade, last_close, update_date").
		From(tableNameTradeStockOrder).
		Where(squirrel.Eq{"order_id": orderID}).
		Join("basic_stock ON trade_stock_order.stock_num = basic_stock.number").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.StockOrder{Stock: new(entity.Stock)}
	if err := row.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Lot, &e.Share,
		&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// QueryAllStockOrderByDate -.
func (r *trade) QueryAllStockOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.StockOrder, error) {
	sql, arg, err := r.Builder.
		Select("order_id, status, order_time, stock_num, action, price, lot, share, number, name, exchange, category, day_trade, last_close, update_date").
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
		if err := rows.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Lot, &e.Share,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryAllStockOrder -.
func (r *trade) QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	sql, _, err := r.Builder.
		Select("order_id, status, order_time, stock_num, action, price, lot, share, number, name, exchange, category, day_trade, last_close, update_date").
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
		if err := rows.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.StockNum, &e.Action, &e.Price, &e.Lot, &e.Share,
			&e.Stock.Number, &e.Stock.Name, &e.Stock.Exchange, &e.Stock.Category, &e.Stock.DayTrade, &e.Stock.LastClose, &e.Stock.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// InsertOrUpdateStockTradeBalance -.
func (r *trade) InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error {
	dbTradeBalance, err := r.queryStockTradeBalanceByDate(ctx, t.TradeDay)
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
		builder := r.Builder.
			Insert(tableNameTradeStockBalance).
			Columns("trade_count, forward, reverse, original_balance, discount, total, trade_day")
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

// queryStockTradeBalanceByDate -.
func (r *trade) queryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.StockTradeBalance, error) {
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
func (r *trade) QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error) {
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

// QueryAllFutureTradeBalance -.
func (r *trade) QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error) {
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

// queryFutureOrderByID -.
func (r *trade) queryFutureOrderByID(ctx context.Context, orderID string) (*entity.FutureOrder, error) {
	sql, arg, err := r.Builder.
		Select("order_id, status, order_time, trade_future_order.code, action, price, position, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
		From(tableNameTradeFutureOrder).
		Where(squirrel.Eq{"order_id": orderID}).
		Join("basic_future ON trade_future_order.code = basic_future.code").ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.FutureOrder{Future: new(entity.Future)}
	if err := row.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.Code, &e.Action, &e.Price, &e.Position,
		&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// InsertOrUpdateFutureOrderByOrderID -.
func (r *trade) InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error {
	dbOrder, err := r.queryFutureOrderByID(ctx, t.OrderID)
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
		builder := r.Builder.
			Insert(tableNameTradeFutureOrder).
			Columns("order_id, status, order_time, code, action, price, position")
		builder = builder.Values(t.OrderID, t.Status, t.OrderTime, t.Code, t.Action, t.Price, t.Position)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else if !cmp.Equal(t, dbOrder) {
		builder := r.Builder.
			Update(tableNameTradeFutureOrder).
			Set("order_id", t.OrderID).
			Set("status", t.Status).
			Set("order_time", t.OrderTime).
			Set("code", t.Code).
			Set("action", t.Action).
			Set("price", t.Price).
			Set("position", t.Position).
			Where(squirrel.Eq{"order_id": t.OrderID})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *trade) QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
	sql, _, err := r.Builder.
		Select("order_id, status, order_time, trade_future_order.code, action, price, position, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
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
		if err := rows.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.Code, &e.Action, &e.Price, &e.Position,
			&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// QueryAllFutureOrderByDate -.
func (r *trade) QueryAllFutureOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.FutureOrder, error) {
	sql, arg, err := r.Builder.
		Select("order_id, status, order_time, trade_future_order.code, action, price, position, basic_future.code, symbol, name, category, delivery_month, delivery_date, underlying_kind, unit, limit_up, limit_down, reference, update_date").
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
		if err := rows.Scan(&e.OrderID, &e.Status, &e.OrderTime, &e.Code, &e.Action, &e.Price, &e.Position,
			&e.Future.Code, &e.Future.Symbol, &e.Future.Name, &e.Future.Category, &e.Future.DeliveryMonth, &e.Future.DeliveryDate, &e.Future.UnderlyingKind, &e.Future.Unit, &e.Future.LimitUp, &e.Future.LimitDown, &e.Future.Reference, &e.Future.UpdateDate); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

// queryFutureTradeBalanceByDate -.
func (r *trade) queryFutureTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.FutureTradeBalance, error) {
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
func (r *trade) InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.FutureTradeBalance) error {
	dbTradeBalance, err := r.queryFutureTradeBalanceByDate(ctx, t.TradeDay)
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
		builder := r.Builder.
			Insert(tableNameFutureTradeBalance).
			Columns("trade_count, forward, reverse, total, trade_day")
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

func (r *trade) QueryLastAccountBalance(ctx context.Context) (*entity.AccountBalance, error) {
	sql, arg, err := r.Builder.
		Select("id, date, balance, today_margin, available_margin, yesterday_margin, risk_indicator").
		From(tableNameAccountBalance).
		Limit(1).OrderBy("date DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.AccountBalance{}
	if err := row.Scan(&e.ID, &e.Date, &e.Balance, &e.TodayMargin, &e.AvailableMargin, &e.YesterdayMargin, &e.RiskIndicator); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *trade) queryAccountBalanceByDate(ctx context.Context, date time.Time) (*entity.AccountBalance, error) {
	sql, arg, err := r.Builder.
		Select("id, date, balance, today_margin, available_margin, yesterday_margin, risk_indicator").
		From(tableNameAccountBalance).
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.AccountBalance{}
	if err := row.Scan(&e.ID, &e.Date, &e.Balance, &e.TodayMargin, &e.AvailableMargin, &e.YesterdayMargin, &e.RiskIndicator); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *trade) InsertOrUpdateAccountBalance(ctx context.Context, t *entity.AccountBalance) error {
	dbStatus, err := r.queryAccountBalanceByDate(ctx, t.Date)
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

	if dbStatus == nil {
		builder := r.Builder.
			Insert(tableNameAccountBalance).
			Columns("date, balance, today_margin, available_margin, yesterday_margin, risk_indicator")
		builder = builder.Values(t.Date, t.Balance, t.TodayMargin, t.AvailableMargin, t.YesterdayMargin, t.RiskIndicator)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameAccountBalance).
			Set("balance", t.Balance).
			Set("today_margin", t.TodayMargin).
			Set("available_margin", t.AvailableMargin).
			Set("yesterday_margin", t.YesterdayMargin).
			Set("risk_indicator", t.RiskIndicator).
			Where(squirrel.Eq{"date": t.Date})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *trade) queryAccountSettlementByDate(ctx context.Context, date time.Time) (*entity.Settlement, error) {
	sql, arg, err := r.Builder.
		Select("date, settlement").
		From(tableNameAccountSettlement).
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.Settlement{}
	if err := row.Scan(&e.Date, &e.Settlement); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *trade) InsertOrUpdateAccountSettlement(ctx context.Context, t *entity.Settlement) error {
	dbSettle, err := r.queryAccountSettlementByDate(ctx, t.Date)
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

	if dbSettle == nil {
		builder := r.Builder.
			Insert(tableNameAccountSettlement).
			Columns("date, settlement")
		builder = builder.Values(t.Date, t.Settlement)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameAccountSettlement).
			Set("settlement", t.Settlement).
			Where(squirrel.Eq{"date": t.Date})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *trade) QueryInventoryUUIDStockByDate(ctx context.Context, date time.Time) (map[string]string, error) {
	sql, arg, err := r.Builder.
		Select("uuid, stock_num").
		From(tableNameInventoryStock).
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := entity.InventoryStock{}
		if err := rows.Scan(
			&e.UUID,
			&e.StockNum,
		); err != nil {
			return nil, err
		}
		result[e.StockNum] = e.UUID
	}
	return result, nil
}

func (r *trade) getInventoryStockByUUID(ctx context.Context, tx pgx.Tx, uuid string) (*entity.InventoryStock, error) {
	sql, args, err := r.Builder.
		Select("uuid, avg_price, lot, share, date, stock_num").
		From(tableNameInventoryStock).
		Where(squirrel.Eq{"uuid": uuid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, sql, args...)
	e := entity.InventoryStock{}
	if err := row.Scan(&e.UUID, &e.AvgPrice, &e.Lot, &e.Share, &e.Date, &e.StockNum); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *trade) insertInventoryStock(ctx context.Context, tx pgx.Tx, t *entity.InventoryStock) error {
	builder := r.Builder.
		Insert(tableNameInventoryStock).
		Columns("uuid, avg_price, lot, share, date, stock_num")
	builder = builder.Values(t.UUID, t.AvgPrice, t.Lot, t.Share, t.Date, t.StockNum)
	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *trade) deleteInventoryStockByUUID(ctx context.Context, tx pgx.Tx, uuid string) error {
	builder := r.Builder.
		Delete(tableNameInventoryStock).
		Where(squirrel.Eq{"uuid": uuid})
	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *trade) getPositionStockByInvID(ctx context.Context, tx pgx.Tx, invID string) ([]*entity.PositionStock, error) {
	sql, args, err := r.Builder.
		Select("inv_id, stock_num, date, quantity, price, last_price, dseq, direction, pnl, fee").
		From(tableNameInventoryPositionStock).
		Where(squirrel.Eq{"inv_id": invID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.PositionStock
	for rows.Next() {
		e := entity.PositionStock{}
		if err = rows.Scan(
			&e.InvID,
			&e.StockNum,
			&e.Date,
			&e.Quantity,
			&e.Price,
			&e.LastPrice,
			&e.Dseq,
			&e.Direction,
			&e.Pnl,
			&e.Fee,
		); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *trade) deletePositionStockByInvID(ctx context.Context, tx pgx.Tx, invID string) error {
	builder := r.Builder.
		Delete(tableNameInventoryPositionStock).
		Where(squirrel.Eq{"inv_id": invID})
	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *trade) insertPositionStock(ctx context.Context, tx pgx.Tx, t []*entity.PositionStock) error {
	builder := r.Builder.
		Insert(tableNameInventoryPositionStock).
		Columns("inv_id, stock_num, date, quantity, price, last_price, dseq, direction, pnl, fee")
	for _, v := range t {
		builder = builder.Values(v.InvID, v.StockNum, v.Date, v.Quantity, v.Price, v.LastPrice, v.Dseq, v.Direction, v.Pnl, v.Fee)
	}
	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *trade) replaceInventoryPositionStock(ctx context.Context, tx pgx.Tx, t []*entity.PositionStock) error {
	if err := r.deletePositionStockByInvID(ctx, tx, t[0].InvID); err != nil {
		return err
	}
	return r.insertPositionStock(ctx, tx, t)
}

func (r *trade) marshalPosition(position []*entity.PositionStock) []byte {
	sort.SliceStable(position, func(i, j int) bool {
		if position[i].Date.Equal(position[j].Date) {
			if position[i].Price == position[j].Price {
				return position[i].Dseq < position[j].Dseq
			}
			return position[i].Price < position[j].Price
		}
		return position[i].Date.Before(position[j].Date)
	})
	b, err := json.Marshal(map[string]any{"position": position})
	if err != nil {
		return nil
	}
	return b
}

func (r *trade) insertOrUpdateInventoryPositionStock(ctx context.Context, tx pgx.Tx, invID string, position []*entity.PositionStock) error {
	if len(position) == 0 {
		return nil
	}

	if dbPosition, err := r.getPositionStockByInvID(ctx, tx, invID); err != nil {
		return err
	} else if len(dbPosition) == 0 {
		if err = r.insertPositionStock(ctx, tx, position); err != nil {
			return err
		}
	} else if !bytes.Equal(r.marshalPosition(position), r.marshalPosition(dbPosition)) {
		if err = r.replaceInventoryPositionStock(ctx, tx, position); err != nil {
			return err
		}
	}
	return nil
}

func (r *trade) InsertOrUpdateInventoryStock(ctx context.Context, t []*entity.InventoryStock) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	for _, v := range t {
		if err = r.insertOrUpdateInventoryPositionStock(ctx, tx, v.UUID, v.Position); err != nil {
			return err
		}
		dbStock, err := r.getInventoryStockByUUID(ctx, tx, v.UUID)
		if err != nil {
			return err
		}
		if dbStock == nil {
			if err = r.insertInventoryStock(ctx, tx, v); err != nil {
				return err
			}
		} else if !cmp.Equal(v, dbStock) {
			if err = r.deleteInventoryStockByUUID(ctx, tx, v.UUID); err != nil {
				return err
			} else if err = r.insertInventoryStock(ctx, tx, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *trade) ClearInventoryStockByUUID(ctx context.Context, uuid string) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	if err := r.deleteInventoryStockByUUID(ctx, tx, uuid); err != nil {
		return err
	}
	if err := r.deletePositionStockByInvID(ctx, tx, uuid); err != nil {
		return err
	}
	return nil
}

func (r *trade) QueryInventoryStockByDate(ctx context.Context, date time.Time) ([]*entity.InventoryStock, error) {
	tx, err := r.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer r.EndTransaction(tx, err)

	sql, arg, err := r.Builder.
		Select("uuid, avg_price, lot, share, date, stock_num").
		From(tableNameInventoryStock).
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	result := []*entity.InventoryStock{}
	rows, err := tx.Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := entity.InventoryStock{}
		if err := rows.Scan(
			&e.UUID,
			&e.AvgPrice,
			&e.Lot,
			&e.Share,
			&e.Date,
			&e.StockNum,
		); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}

	for i, v := range result {
		position, err := r.getPositionStockByInvID(ctx, tx, v.UUID)
		if err != nil {
			continue
		}
		result[i].Position = position
	}
	return result, nil
}
