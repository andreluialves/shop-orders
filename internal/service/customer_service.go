package service

import (
	"context"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
)

type CustomerService struct {
	customerRepository repository.CustomerRepository
	idGenerator        repository.CustomerIDGenerator
}

func NewCustomerService(customerRepository repository.CustomerRepository, idGenerator repository.CustomerIDGenerator) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		idGenerator:        idGenerator,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	customer.ID, _ = s.idGenerator.NextCustomerID(ctx)
	return s.customerRepository.Save(customer)
}

func (s *CustomerService) FindByID(id string) (*domain.Customer, error) {
	return s.customerRepository.FindByID(id)
}

func (s *CustomerService) List() ([]*domain.Customer, error) {
	return s.customerRepository.List()
}
