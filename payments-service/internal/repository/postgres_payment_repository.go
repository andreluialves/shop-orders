package repository

import (
	"context"
	"errors"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
	"github.com/jackc/pgx/v5"
)

type PostgresPaymentRepository struct {
	db DBTX
}

func NewPostgresPaymentRepository(db DBTX) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{
		db: db,
	}
}

func (r *PostgresPaymentRepository) Save(payment *domain.Payment) error {

	query := `
		INSERT INTO payments (
			id,
			order_id,
			amount,
			status
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		payment.ID,
		payment.OrderID,
		payment.Amount,
		string(payment.Status()),
	)

	return err
}

func (r *PostgresPaymentRepository) FindByID(id string) (*domain.Payment, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			order_id,
			amount,
			status
		FROM payments
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var (
		paymentID string
		orderID   string
		amount    float64
		status    string
	)

	err := row.Scan(
		&paymentID,
		&orderID,
		&amount,
		&status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	payment := domain.RestorePayment(
		paymentID,
		orderID,
		amount,
		domain.PaymentStatus(status),
	)

	return payment, nil
}

func (r *PostgresPaymentRepository) List() ([]*domain.Payment, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			order_id,
			amount,
			status
		FROM payments
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment

	for rows.Next() {

		var (
			id      string
			orderID string
			amount  float64
			status  string
		)

		if err := rows.Scan(
			&id,
			&orderID,
			&amount,
			&status,
		); err != nil {
			return nil, err
		}

		payment := domain.RestorePayment(
			id,
			orderID,
			amount,
			domain.PaymentStatus(status),
		)

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
