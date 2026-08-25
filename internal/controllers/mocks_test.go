package controllers_test

import (
	"context"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
	"github.com/andreluialves/shop-orders/internal/service"
)

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

type mockUnitOfWork struct {
	repos repository.Repositories
}

func (m *mockUnitOfWork) Execute(ctx context.Context, fn func(repos repository.Repositories) error) error {
	return fn(m.repos)
}

func newTestOrderService(productRepo *mockProductRepository, orderRepo *mockOrderRepository) *service.OrderService {
	uow := &mockUnitOfWork{
		repos: repository.Repositories{
			Order:   orderRepo,
			Product: productRepo,
		},
	}

	return service.NewOrderService(productRepo, orderRepo, uow)
}
