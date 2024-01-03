// Package repo package repo
package repo

import (
	"context"
	"errors"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4"
)

type TradeRepo struct {
	*postgres.Postgres
}

func NewTrade(pg *postgres.Postgres) *TradeRepo {
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
		builder := r.Builder.Insert(tableNameTradeStockOrder).Columns("order_id, status, order_time, stock_num, action, price, lot, share")
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

// QueryStockOrderByID -.
func (r *TradeRepo) QueryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error) {
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
func (r *TradeRepo) QueryAllStockOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.StockOrder, error) {
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
func (r *TradeRepo) QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
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
		builder := r.Builder.Insert(tableNameTradeFutureOrder).Columns("order_id, status, order_time, code, action, price, position")
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

func (r *TradeRepo) QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
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
func (r *TradeRepo) QueryAllFutureOrderByDate(ctx context.Context, timeRange []time.Time) ([]*entity.FutureOrder, error) {
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

func (r *TradeRepo) QueryAllLastAccountBalance(ctx context.Context, bankIDArr []int) ([]*entity.AccountBalance, error) {
	var result []*entity.AccountBalance
	for _, v := range bankIDArr {
		sql, arg, err := r.Builder.
			Select("id, date, balance, today_margin, available_margin, yesterday_margin, risk_indicator, bank_id").
			From(tableNameAccountBalance).
			Where(squirrel.Eq{"bank_id": v}).
			Limit(1).OrderBy("date DESC").
			ToSql()
		if err != nil {
			return nil, err
		}

		row := r.Pool().QueryRow(ctx, sql, arg...)
		e := entity.AccountBalance{}
		if err := row.Scan(&e.ID, &e.Date, &e.Balance, &e.TodayMargin, &e.AvailableMargin, &e.YesterdayMargin, &e.RiskIndicator, &e.BankID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *TradeRepo) QueryAccountBalanceByDateAndBankID(ctx context.Context, date time.Time, bankID int) (*entity.AccountBalance, error) {
	sql, arg, err := r.Builder.
		Select("id, date, balance, today_margin, available_margin, yesterday_margin, risk_indicator, bank_id").
		From(tableNameAccountBalance).
		Where(squirrel.Eq{"date": date, "bank_id": bankID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.AccountBalance{}
	if err := row.Scan(&e.ID, &e.Date, &e.Balance, &e.TodayMargin, &e.AvailableMargin, &e.YesterdayMargin, &e.RiskIndicator, &e.BankID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *TradeRepo) InsertOrUpdateAccountBalance(ctx context.Context, t *entity.AccountBalance) error {
	dbStatus, err := r.QueryAccountBalanceByDateAndBankID(ctx, t.Date, t.BankID)
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
		builder := r.Builder.Insert(tableNameAccountBalance).Columns("date, balance, today_margin, available_margin, yesterday_margin, risk_indicator, bank_id")
		builder = builder.Values(t.Date, t.Balance, t.TodayMargin, t.AvailableMargin, t.YesterdayMargin, t.RiskIndicator, t.BankID)
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
			Where(squirrel.Eq{"date": t.Date, "bank_id": t.BankID})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *TradeRepo) QueryAccountSettlementByDate(ctx context.Context, date time.Time) (*entity.Settlement, error) {
	sql, arg, err := r.Builder.
		Select("date, sinopac, fugle").
		From(tableNameAccountSettlement).
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.Settlement{}
	if err := row.Scan(&e.Date, &e.Sinopac, &e.Fugle); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *TradeRepo) InsertOrUpdateAccountSettlement(ctx context.Context, t *entity.Settlement) error {
	dbSettle, err := r.QueryAccountSettlementByDate(ctx, t.Date)
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
		builder := r.Builder.Insert(tableNameAccountSettlement).Columns("date, sinopac, fugle")
		builder = builder.Values(t.Date, t.Sinopac, t.Fugle)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	} else {
		builder := r.Builder.
			Update(tableNameAccountSettlement).
			Set("sinopac", t.Sinopac).
			Set("fugle", t.Fugle).
			Where(squirrel.Eq{"date": t.Date})
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *TradeRepo) QueryInventoryStockByDate(ctx context.Context, date time.Time) ([]*entity.InventoryStock, error) {
	sql, arg, err := r.Builder.
		Select("id, bank_id, avg_price, lot, share, updated, stock_num").
		From(tableNameInventoryStock).
		Where(squirrel.Eq{"updated": date}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.InventoryStock
	for rows.Next() {
		e := entity.InventoryStock{}
		if err := rows.Scan(
			&e.ID,
			&e.BankID,
			&e.AvgPrice,
			&e.Lot,
			&e.Share,
			&e.Updated,
			&e.StockNum,
		); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *TradeRepo) DeleteInventoryStockByDate(ctx context.Context, date time.Time) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	builder := r.Builder.Delete(tableNameInventoryStock).
		Where(squirrel.Eq{"updated": date})
	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *TradeRepo) InsertInventoryStock(ctx context.Context, t []*entity.InventoryStock) error {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	for _, v := range t {
		builder := r.Builder.Insert(tableNameInventoryStock).Columns("bank_id, avg_price, lot, share, updated, stock_num")
		builder = builder.Values(v.BankID, v.AvgPrice, v.Lot, v.Share, v.Updated, v.StockNum)
		if sql, args, err = builder.ToSql(); err != nil {
			return err
		} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}
