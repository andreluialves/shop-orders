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
	customerRepo := &mockCustomerRepository{
		FindByIDFunc: func(id string) (*domain.Customer, error) {
			return &domain.Customer{ID: id, Name: "Cliente Teste"}, nil
		},
	}

	return newTestOrderServiceWithCustomerRepo(productRepo, orderRepo, customerRepo)
}

// newTestOrderServiceWithCustomerRepo permite controlar o comportamento do
// CustomerRepository nos testes que precisam simular cliente não encontrado.
func newTestOrderServiceWithCustomerRepo(
	productRepo *mockProductRepository,
	orderRepo *mockOrderRepository,
	customerRepo *mockCustomerRepository,
) *service.OrderService {
	uow := &mockUnitOfWork{
		repos: repository.Repositories{
			Order:    orderRepo,
			Product:  productRepo,
			Customer: customerRepo,
		},
	}

	idGenerator := &mockOrderIDGenerator{}

	return service.NewOrderService(productRepo, orderRepo, uow, idGenerator)
}

type mockOrderIDGenerator struct {
	NextOrderIDFunc func(ctx context.Context) (string, error)
}

func (m *mockOrderIDGenerator) NextOrderID(ctx context.Context) (string, error) {
	if m.NextOrderIDFunc != nil {
		return m.NextOrderIDFunc(ctx)
	}
	return "PED-TEST", nil
}

type mockCustomerRepository struct {
	FindByIDFunc func(id string) (*domain.Customer, error)
	SaveFunc     func(customer *domain.Customer) error
	ListFunc     func() ([]*domain.Customer, error)
}

func (m *mockCustomerRepository) FindByID(id string) (*domain.Customer, error) {
	return m.FindByIDFunc(id)
}

func (m *mockCustomerRepository) Save(customer *domain.Customer) error {
	return m.SaveFunc(customer)
}

func (m *mockCustomerRepository) List() ([]*domain.Customer, error) {
	return m.ListFunc()
}

type mockCustomerIDGenerator struct {
	NextCustomerIDFunc func(ctx context.Context) (string, error)
}

func (m *mockCustomerIDGenerator) NextCustomerID(ctx context.Context) (string, error) {
	if m.NextCustomerIDFunc != nil {
		return m.NextCustomerIDFunc(ctx)
	}
	return "CUST-TEST", nil
}
