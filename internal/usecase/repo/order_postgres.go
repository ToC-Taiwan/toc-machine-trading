// Package repo package repo
package repo

import (
	"context"
	"errors"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"

	"github.com/Masterminds/squirrel"
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

// InserOrUpdateTradeBalance -.
func (r *OrderRepo) InserOrUpdateTradeBalance(ctx context.Context, t *entity.TradeBalance) error {
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
