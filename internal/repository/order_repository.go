package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	List() ([]*domain.Order, error)
}
