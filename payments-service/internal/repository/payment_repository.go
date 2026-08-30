package repository

import (
	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByID(id string) (*domain.Payment, error)
	List() ([]*domain.Payment, error)
}
