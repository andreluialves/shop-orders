package repository

import (
	"context"
	"errors"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (r *PostgresProductRepository) Save(product *domain.Product) error {

	query := `
		INSERT INTO products (
			id,
			name,
			price,
			quantity
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			price = EXCLUDED.price,
			quantity = EXCLUDED.quantity
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Quantity,
	)

	return err
}

func (r *PostgresProductRepository) FindByID(id string) (*domain.Product, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			name,
			price,
			quantity
		FROM products
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var (
		productID string
		name      string
		price     float64
		quantity  int
	)

	if err := row.Scan(
		&productID,
		&name,
		&price,
		&quantity,
	); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}

		return nil, err
	}

	product := domain.RestoreProduct(
		productID,
		name,
		price,
		quantity,
	)

	return product, nil
}

func (r *PostgresProductRepository) List() ([]*domain.Product, error) {

	query := `
		SELECT
			id,
			name,
			price,
			quantity
		FROM products
		ORDER BY name
	`

	rows, err := r.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*domain.Product

	for rows.Next() {

		var product domain.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
