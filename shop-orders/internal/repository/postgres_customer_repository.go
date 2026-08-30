package repository

import (
	"context"
	"errors"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/jackc/pgx/v5"
)

type PostgresCustomerRepository struct {
	db DBTX
}

func NewPostgresCustomerRepository(db DBTX) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{
		db: db,
	}
}

func (r *PostgresCustomerRepository) Save(customer *domain.Customer) error {

	query := `
		INSERT INTO customers (
			id,
			name,
			email,
			address,
			phone
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Address,
		customer.Phone,
	)

	return err
}

func (r *PostgresCustomerRepository) FindByID(id string) (*domain.Customer, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			name,
			email,
			address,
			phone
		FROM customers
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	customer := &domain.Customer{}

	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Address,
		&customer.Phone,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCustomerNotFound
	}

	return customer, nil
}

func (r *PostgresCustomerRepository) List() ([]*domain.Customer, error) {

	query := `
		SELECT
			id,
			name,
			email,
			address,
			phone
		FROM customers
		ORDER BY name
	`

	rows, err := r.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*domain.Customer

	for rows.Next() {

		var customer domain.Customer

		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Address,
			&customer.Phone,
		)

		if err != nil {
			return nil, err
		}

		customers = append(customers, &customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
