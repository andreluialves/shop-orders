package service_test

import (
	"context"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
)

type mockPaymentRepository struct {
	SaveFunc     func(payment *domain.Payment) error
	FindByIDFunc func(id string) (*domain.Payment, error)
	ListFunc     func() ([]*domain.Payment, error)
}

func (m *mockPaymentRepository) Save(payment *domain.Payment) error {
	return m.SaveFunc(payment)
}

func (m *mockPaymentRepository) FindByID(id string) (*domain.Payment, error) {
	return m.FindByIDFunc(id)
}

func (m *mockPaymentRepository) List() ([]*domain.Payment, error) {
	return m.ListFunc()
}

type mockPaymentIDGenerator struct {
	NextPaymentIDFunc func(ctx context.Context) (string, error)
}

func (m *mockPaymentIDGenerator) NextPaymentID(ctx context.Context) (string, error) {
	if m.NextPaymentIDFunc != nil {
		return m.NextPaymentIDFunc(ctx)
	}
	return "PAY-TEST", nil
}
