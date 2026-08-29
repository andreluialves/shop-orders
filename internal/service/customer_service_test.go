package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/service"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
	t.Run("deve criar cliente e atribuir ID gerado", func(t *testing.T) {
		var savedCustomer *domain.Customer

		customerRepo := &mockCustomerRepository{
			SaveFunc: func(c *domain.Customer) error {
				savedCustomer = c
				return nil
			},
		}

		idGenerator := &mockCustomerIDGenerator{
			NextCustomerIDFunc: func(ctx context.Context) (string, error) {
				return "CUST-001", nil
			},
		}

		s := service.NewCustomerService(customerRepo, idGenerator)

		customer := &domain.Customer{
			Name:    "João Silva",
			Email:   "joao@example.com",
			Address: "Rua A, 123",
			Phone:   "11999999999",
		}

		err := s.CreateCustomer(context.Background(), customer)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if customer.ID != "CUST-001" {
			t.Errorf("esperava ID CUST-001, recebeu %v", customer.ID)
		}

		if savedCustomer == nil {
			t.Fatal("esperava que Save fosse chamado")
		}
	})

	t.Run("deve propagar erro quando geração de ID falha", func(t *testing.T) {
		var saveCalled bool

		customerRepo := &mockCustomerRepository{
			SaveFunc: func(c *domain.Customer) error {
				saveCalled = true
				return nil
			},
		}

		idGenerator := &mockCustomerIDGenerator{
			NextCustomerIDFunc: func(ctx context.Context) (string, error) {
				return "", errors.New("falha ao gerar ID")
			},
		}

		s := service.NewCustomerService(customerRepo, idGenerator)

		customer := &domain.Customer{Name: "João Silva", Email: "joao@example.com", Address: "Rua A, 123"}

		err := s.CreateCustomer(context.Background(), customer)

		if err == nil {
			t.Error("esperava erro propagado do gerador de ID")
		}

		if saveCalled {
			t.Error("Save não deveria ser chamado quando a geração de ID falha")
		}
	})

	t.Run("deve propagar erro quando Save falha", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			SaveFunc: func(c *domain.Customer) error {
				return errors.New("erro de conexão com banco")
			},
		}

		s := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})

		customer := &domain.Customer{Name: "João Silva", Email: "joao@example.com", Address: "Rua A, 123"}

		err := s.CreateCustomer(context.Background(), customer)

		if err == nil {
			t.Error("esperava erro propagado do Save")
		}
	})
}

func TestCustomerService_FindByID(t *testing.T) {
	t.Run("deve retornar cliente quando encontrado", func(t *testing.T) {
		expected := &domain.Customer{ID: "CUST-001", Name: "João Silva", Email: "joao@example.com"}

		customerRepo := &mockCustomerRepository{
			FindByIDFunc: func(id string) (*domain.Customer, error) {
				return expected, nil
			},
		}

		s := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})

		result, err := s.FindByID("CUST-001")

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if result != expected {
			t.Error("esperava o cliente retornado pelo repository")
		}
	})

	t.Run("deve retornar erro quando cliente não é encontrado", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			FindByIDFunc: func(id string) (*domain.Customer, error) {
				return nil, domain.ErrCustomerNotFound
			},
		}

		s := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})

		_, err := s.FindByID("CUST-999")

		if !errors.Is(err, domain.ErrCustomerNotFound) {
			t.Errorf("esperava ErrCustomerNotFound, recebeu %v", err)
		}
	})
}

func TestCustomerService_List(t *testing.T) {
	t.Run("deve retornar lista de clientes", func(t *testing.T) {
		expected := []*domain.Customer{
			{ID: "CUST-001", Name: "João Silva"},
			{ID: "CUST-002", Name: "Maria Souza"},
		}

		customerRepo := &mockCustomerRepository{
			ListFunc: func() ([]*domain.Customer, error) {
				return expected, nil
			},
		}

		s := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})

		result, err := s.List()

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if len(result) != 2 {
			t.Errorf("esperava 2 clientes, recebeu %d", len(result))
		}
	})

	t.Run("deve propagar erro quando List falha", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			ListFunc: func() ([]*domain.Customer, error) {
				return nil, errors.New("erro de conexão com banco")
			},
		}

		s := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})

		_, err := s.List()

		if err == nil {
			t.Error("esperava erro propagado do List")
		}
	})
}
