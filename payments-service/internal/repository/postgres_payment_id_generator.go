package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPaymentIDGenerator struct {
	pool *pgxpool.Pool
}

func NewPostgresPaymentIDGenerator(pool *pgxpool.Pool) *PostgresPaymentIDGenerator {
	return &PostgresPaymentIDGenerator{pool: pool}
}

func (g *PostgresPaymentIDGenerator) NextPaymentID(ctx context.Context) (string, error) {
	var next int64

	err := g.pool.QueryRow(ctx, `SELECT nextval('payments_id_seq')`).Scan(&next)
	if err != nil {
		return "", fmt.Errorf("gerar próximo ID de pagamento: %w", err)
	}

	return fmt.Sprintf("PAY-%03d", next), nil
}
