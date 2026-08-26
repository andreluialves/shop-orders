package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCustomerIDGenerator struct {
	pool *pgxpool.Pool
}

func NewPostgresCustomerIDGenerator(pool *pgxpool.Pool) *PostgresCustomerIDGenerator {
	return &PostgresCustomerIDGenerator{pool: pool}
}

func (g *PostgresCustomerIDGenerator) NextCustomerID(ctx context.Context) (string, error) {
	var next int64

	err := g.pool.QueryRow(ctx, `SELECT nextval('customers_id_seq')`).Scan(&next)
	if err != nil {
		return "", fmt.Errorf("gerar próximo ID de cliente: %w", err)
	}

	return fmt.Sprintf("CUST-%03d", next), nil
}
