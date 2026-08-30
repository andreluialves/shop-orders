package repository

import (
	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	Update(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	List(limit, offset int) ([]*domain.Order, int, error)
}
