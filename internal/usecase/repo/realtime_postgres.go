package repo

import (
	"context"

	"tmt/internal/entity"
	"tmt/pkg/postgres"
)

// RealTimeRepo -.
type RealTimeRepo struct {
	*postgres.Postgres
}

func NewRealTime(pg *postgres.Postgres) *RealTimeRepo {
	return &RealTimeRepo{pg}
}

// InsertEvent -.
func (r *RealTimeRepo) InsertEvent(ctx context.Context, t *entity.SinopacEvent) error {
	builder := r.Builder.Insert(tableNameEvent).
		Columns("event, event_code, info, response, event_time").
		Values(t.Event, t.EventCode, t.Info, t.Response, t.EventTime)

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
