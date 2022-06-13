package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
)

// StreamRepo -.
type StreamRepo struct {
	*postgres.Postgres
}

// NewStream -.
func NewStream(pg *postgres.Postgres) *StreamRepo {
	return &StreamRepo{pg}
}

// InsertEvent -.
func (r *StreamRepo) InsertEvent(ctx context.Context, t *entity.SinopacEvent) error {
	builder := r.Builder.Insert(tableNameEvent).
		Columns("event, event_code, info, response, event_time").
		Values(t.Event, t.EventCode, t.Info, t.Response, t.EventTime)

	if sql, args, err := builder.ToSql(); err != nil {
		return err
	} else if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
