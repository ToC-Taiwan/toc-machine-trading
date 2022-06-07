package repo

import (
	"context"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/postgres"
)

const (
	tableNameEvent string = "sinopac_event"
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
		Columns("event, event_code, info, response").
		Values(t.Event, t.EventCode, t.Info, t.Response)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
