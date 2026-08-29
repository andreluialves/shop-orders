package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
)

type CustomerRepository interface {
	Save(customer *domain.Customer) error
	FindByID(id string) (*domain.Customer, error)
	List() ([]*domain.Customer, error)
}
