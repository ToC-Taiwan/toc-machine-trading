package repo

import (
	"context"
	"errors"

	"tmt/internal/entity"
	"tmt/pkg/postgres"

	"github.com/jackc/pgx/v4"
)

// SystemRepo -.
type SystemRepo struct {
	*postgres.Postgres
}

func NewSystemRepo(pg *postgres.Postgres) *SystemRepo {
	return &SystemRepo{pg}
}

func (r *SystemRepo) InsertUser(ctx context.Context, t *entity.User) error {
	builder := r.Builder.Insert(tableNameSystemAccount).
		Columns("username, password, email").
		Values(t.Username, t.Password, t.Email)

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *SystemRepo) EmailVerification(ctx context.Context, username string) error {
	builder := r.Builder.Update(tableNameSystemAccount).
		Set("email_verified", true).
		Where("username = ?", username)

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}
	defer r.EndTransaction(tx, err)
	var sql string
	var args []interface{}

	if sql, args, err = builder.ToSql(); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *SystemRepo) QueryUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	sql, arg, err := r.Builder.
		Select("username, password, email, email_verified").
		From(tableNameSystemAccount).
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.User{}
	if err := row.Scan(&e.Username, &e.Password, &e.Email, &e.EmailVerified); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *SystemRepo) QueryAllUser(ctx context.Context) ([]*entity.User, error) {
	sql, arg, err := r.Builder.
		Select("username, email").
		From(tableNameSystemAccount).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.User
	for rows.Next() {
		e := entity.User{}
		if err := rows.Scan(&e.Username, &e.Email); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}
