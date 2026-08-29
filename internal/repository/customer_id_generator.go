package repository

import "context"

type CustomerIDGenerator interface {
	NextCustomerID(ctx context.Context) (string, error)
}
