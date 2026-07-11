package repository

import (
	"fmt"

	"github.com/andreluialves/shop-orders/internal/domain"
)

type MemoryOrderRepository struct {
	orders map[string]*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *MemoryOrderRepository) Save(order *domain.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) FindByID(id string) (*domain.Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order %s: %w", id, domain.ErrOrderNotFound)
	}

	return order, nil
}

func (r *MemoryOrderRepository) List() ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, len(r.orders))

	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}
