package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repositories agrupa os repositories vinculados a uma mesma transação.
type Repositories struct {
	Order    OrderRepository
	Product  ProductRepository
	Customer CustomerRepository
}

// UnitOfWork executa fn dentro de uma transação: se fn retornar erro, tudo é
// desfeito (rollback); se retornar nil, tudo é persistido (commit).
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(repos Repositories) error) error
}

type PostgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{pool: pool}
}

func (u *PostgresUnitOfWork) Execute(ctx context.Context, fn func(repos Repositories) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	repos := Repositories{
		Order:    NewPostgresOrderRepository(tx),
		Product:  NewPostgresProductRepository(tx),
		Customer: NewPostgresCustomerRepository(tx), // novo
	}

	if err := fn(repos); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
