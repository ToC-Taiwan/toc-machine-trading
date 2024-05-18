package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/pkg/postgres"
)

// system -.
type system struct {
	*postgres.Postgres
}

func NewSystemRepo(pg *postgres.Postgres) SystemRepo {
	return &system{pg}
}

func (r *system) InsertUser(ctx context.Context, t *entity.NewUser) error {
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

func (r *system) EmailVerification(ctx context.Context, username string) error {
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

func (r *system) QueryUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	sql, arg, err := r.Builder.
		Select("username, password, email, email_verified, auth_trade").
		From(tableNameSystemAccount).
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.User{}
	if err := row.Scan(&e.Username, &e.Password, &e.Email, &e.EmailVerified, &e.AuthTrade); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *system) queryUserIDByUsername(ctx context.Context, username string) (int, error) {
	sql, arg, err := r.Builder.
		Select("id").
		From(tableNameSystemAccount).
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return 0, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (r *system) QueryAllUser(ctx context.Context) ([]*entity.User, error) {
	sql, arg, err := r.Builder.
		Select("username, email, email_verified, auth_trade").
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
		if err := rows.Scan(&e.Username, &e.Email, &e.EmailVerified, &e.AuthTrade); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	return result, nil
}

func (r *system) InsertOrUpdatePushToken(ctx context.Context, token, username string, enabled bool) error {
	dbToken, err := r.GetPushToken(ctx, token)
	if err != nil {
		return err
	} else if dbToken != nil {
		return r.updatePushToken(ctx, token, enabled)
	}

	userID, err := r.queryUserIDByUsername(ctx, username)
	if err != nil {
		return err
	} else if userID == 0 {
		return errors.New("user not found")
	}

	builder := r.Builder.Insert(tableNameSystemPushToken).
		Columns("created, token, user_id, enabled").
		Values(time.Now(), token, userID, enabled)

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

func (r *system) updatePushToken(ctx context.Context, token string, enabled bool) error {
	builder := r.Builder.Update(tableNameSystemPushToken).
		Set("created", time.Now()).
		Set("enabled", enabled).
		Where("token = ?", token)

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

func (r *system) GetPushToken(ctx context.Context, token string) (*entity.PushToken, error) {
	sql, arg, err := r.Builder.
		Select("created, token, user_id, enabled").
		From(tableNameSystemPushToken).
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Pool().QueryRow(ctx, sql, arg...)
	e := entity.PushToken{}
	if err := row.Scan(&e.Created, &e.Token, &e.UserID, &e.Enabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *system) GetAllPushTokens(ctx context.Context) ([]string, error) {
	sql, arg, err := r.Builder.
		Select("token").
		From(tableNameSystemPushToken).
		Where("enabled = ?", true).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool().Query(ctx, sql, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		if token == "" {
			continue
		}
		result = append(result, token)
	}
	return result, nil
}

func (r *system) DeleteAllPushTokens(ctx context.Context) error {
	builder := r.Builder.Delete(tableNameSystemPushToken)

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

func (r *system) GetLastJWT(ctx context.Context) (string, error) {
	sql, arg, err := r.Builder.
		Select("key").
		From(tableNameSystemJWT).
		OrderBy("id DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return "", err
	}

	rows := r.Pool().QueryRow(ctx, sql, arg...)
	var result string
	if err := rows.Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *system) InsertJWT(ctx context.Context, jwt string) error {
	builder := r.Builder.Insert(tableNameSystemJWT).
		Columns("key, created").
		Values(jwt, time.Now())

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
