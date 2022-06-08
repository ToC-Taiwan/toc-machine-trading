package repo

import "toc-machine-trading/pkg/postgres"

// OrderRepo -.
type OrderRepo struct {
	*postgres.Postgres
}

// NewOrder -.
func NewOrder(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}
