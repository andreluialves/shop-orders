package repository_test

import "github.com/andreluialves/shop-orders/shop-orders/internal/domain"

type mockOrderRepository struct {
	FindByIDFunc func(id string) (*domain.Order, error)
	SaveFunc     func(order *domain.Order) error
	UpdateFunc   func(order *domain.Order) error
	ListFunc     func(limit, offset int) ([]*domain.Order, int, error)
}

func (m *mockOrderRepository) FindByID(id string) (*domain.Order, error) {
	return m.FindByIDFunc(id)
}

func (m *mockOrderRepository) Save(order *domain.Order) error {
	return m.SaveFunc(order)
}

func (m *mockOrderRepository) Update(order *domain.Order) error {
	return m.UpdateFunc(order)
}

func (m *mockOrderRepository) List(limit, offset int) ([]*domain.Order, int, error) {
	return m.ListFunc(limit, offset)
}

type mockProductRepository struct {
	FindByIDFunc func(id string) (*domain.Product, error)
	SaveFunc     func(product *domain.Product) error
	ListFunc     func() ([]*domain.Product, error)
}

func (m *mockProductRepository) FindByID(id string) (*domain.Product, error) {
	return m.FindByIDFunc(id)
}

func (m *mockProductRepository) Save(product *domain.Product) error {
	return m.SaveFunc(product)
}

func (m *mockProductRepository) List() ([]*domain.Product, error) {
	return m.ListFunc()
}

type mockLogger struct {
	InfoFunc  func(msg string, args ...any)
	WarnFunc  func(msg string, args ...any)
	ErrorFunc func(msg string, args ...any)
	DebugFunc func(msg string, args ...any)
}

func (m *mockLogger) Info(msg string, args ...any) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, args...)
	}
}

func (m *mockLogger) Warn(msg string, args ...any) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg, args...)
	}
}

func (m *mockLogger) Error(msg string, args ...any) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, args...)
	}
}

func (m *mockLogger) Debug(msg string, args ...any) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg, args...)
	}
}
