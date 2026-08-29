package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
)

type ProductRepository interface {
	Save(product *domain.Product) error
	FindByID(id string) (*domain.Product, error)
	List() ([]*domain.Product, error)
}
