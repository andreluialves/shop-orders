package repository

import "context"

type OrderIDGenerator interface {
	NextOrderID(ctx context.Context) (string, error)
}
