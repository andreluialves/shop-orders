package repository

import "context"

type PaymentIDGenerator interface {
	NextPaymentID(ctx context.Context) (string, error)
}
