package repository

import (
	"context"
	"errors"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(db *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

func (r *PostgresOrderRepository) Save(order *domain.Order) error {

	ctx := context.Background()

	query := `
		INSERT INTO orders (
			id,
			customer,
			status
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		order.ID,
		order.Customer,
		order.Status(),
	)

	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO order_items (
			order_id,
			product_id,
			quantity,
			price
		)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {

		_, err := r.db.Exec(
			ctx,
			itemQuery,
			order.ID,
			item.Product.ID,
			item.Quantity,
			item.Price,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresOrderRepository) Update(order *domain.Order) error {

	query := `
		UPDATE orders
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		order.Status(),
		order.ID,
	)

	return err
}

func (r *PostgresOrderRepository) FindByID(id string) (*domain.Order, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			customer,
			status
		FROM orders
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var (
		orderID  string
		customer string
		status   string
	)

	if err := row.Scan(
		&orderID,
		&customer,
		&status,
	); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}

		return nil, err
	}

	order := domain.RestoreOrder(
		orderID,
		customer,
		domain.OrderStatus(status),
	)

	itemsQuery := `
		SELECT
			p.id,
			p.name,
			p.price,
			p.quantity,
			oi.quantity,
			oi.price
		FROM order_items oi
		INNER JOIN products p
			ON p.id = oi.product_id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var (
			productID       string
			productName     string
			productPrice    float64
			productQuantity int

			itemQuantity int
			itemPrice    float64
		)

		if err := rows.Scan(
			&productID,
			&productName,
			&productPrice,
			&productQuantity,
			&itemQuantity,
			&itemPrice,
		); err != nil {
			return nil, err
		}

		product := &domain.Product{
			ID:       productID,
			Name:     productName,
			Price:    productPrice,
			Quantity: productQuantity,
		}

		item := domain.NewOrderItem(
			product,
			itemQuantity,
			itemPrice,
		)

		order.AddItem(item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *PostgresOrderRepository) List() ([]*domain.Order, error) {

	ctx := context.Background()

	query := `
		SELECT id
		FROM orders
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {

		var id string

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		order, err := r.FindByID(id)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
