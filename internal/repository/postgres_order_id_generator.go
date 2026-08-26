package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderIDGenerator struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderIDGenerator(pool *pgxpool.Pool) *PostgresOrderIDGenerator {
	return &PostgresOrderIDGenerator{pool: pool}
}

func (g *PostgresOrderIDGenerator) NextOrderID(ctx context.Context) (string, error) {
	var next int64

	err := g.pool.QueryRow(ctx, `SELECT nextval('orders_id_seq')`).Scan(&next)
	if err != nil {
		return "", fmt.Errorf("gerar próximo ID de pedido: %w", err)
	}

	return fmt.Sprintf("PED-%03d", next), nil
}
