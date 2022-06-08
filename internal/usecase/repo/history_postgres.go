package repo

import "toc-machine-trading/pkg/postgres"

// HistoryRepo -.
type HistoryRepo struct {
	*postgres.Postgres
}

// NewHistory -.
func NewHistory(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg}
}
