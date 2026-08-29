package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreluialves/shop-orders/internal/controllers"
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/service"
)

func TestCustomerController_CreateCustomer(t *testing.T) {
	t.Run("deve criar cliente com sucesso e retornar 201", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			SaveFunc: func(c *domain.Customer) error { return nil },
		}
		idGenerator := &mockCustomerIDGenerator{
			NextCustomerIDFunc: func(ctx context.Context) (string, error) {
				return "CUST-001", nil
			},
		}

		customerService := service.NewCustomerService(customerRepo, idGenerator)
		controller := controllers.NewCustomerController(customerService)

		body := `{"name":"João Silva","email":"joao@example.com","address":"Rua A, 123","phone":"11999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateCustomer(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("esperava status 201, recebeu %d", resp.StatusCode)
		}

		var got domain.Customer
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.ID != "CUST-001" {
			t.Errorf("esperava ID CUST-001, recebeu %v", got.ID)
		}
	})

	t.Run("deve retornar 400 quando body é JSON inválido", func(t *testing.T) {
		customerService := service.NewCustomerService(&mockCustomerRepository{}, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		body := `{"name": "João",` // malformado
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateCustomer(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", w.Result().StatusCode)
		}
	})

	t.Run("deve retornar 400 quando email é inválido", func(t *testing.T) {
		customerService := service.NewCustomerService(&mockCustomerRepository{}, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		body := `{"name":"João Silva","email":"marimbinha","address":"Rua A, 123","phone":"11999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateCustomer(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", w.Result().StatusCode)
		}
	})

	t.Run("deve retornar 400 quando nome é vazio", func(t *testing.T) {
		customerService := service.NewCustomerService(&mockCustomerRepository{}, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		body := `{"name":"","email":"joao@example.com","address":"Rua A, 123","phone":"11999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateCustomer(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", w.Result().StatusCode)
		}
	})

	t.Run("deve retornar 500 quando Save falha", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			SaveFunc: func(c *domain.Customer) error {
				return errors.New("erro de conexão com banco")
			},
		}

		customerService := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		body := `{"name":"João Silva","email":"joao@example.com","address":"Rua A, 123","phone":"11999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateCustomer(w, req)

		if w.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", w.Result().StatusCode)
		}
	})
}

func TestCustomerController_FindAllCustomers(t *testing.T) {
	t.Run("deve retornar lista de clientes com 200", func(t *testing.T) {
		expected := []*domain.Customer{
			{ID: "CUST-001", Name: "João Silva"},
			{ID: "CUST-002", Name: "Maria Souza"},
		}

		customerRepo := &mockCustomerRepository{
			ListFunc: func() ([]*domain.Customer, error) {
				return expected, nil
			},
		}

		customerService := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		w := httptest.NewRecorder()

		controller.FindAllCustomers(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var got []*domain.Customer
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if len(got) != 2 {
			t.Errorf("esperava 2 clientes, recebeu %d", len(got))
		}
	})

	t.Run("deve retornar 500 quando List falha", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			ListFunc: func() ([]*domain.Customer, error) {
				return nil, errors.New("erro de conexão com banco")
			},
		}

		customerService := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		w := httptest.NewRecorder()

		controller.FindAllCustomers(w, req)

		if w.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", w.Result().StatusCode)
		}
	})
}

func TestCustomerController_FindCustomerByID(t *testing.T) {
	t.Run("deve retornar cliente com 200 quando encontrado", func(t *testing.T) {
		expected := &domain.Customer{ID: "CUST-001", Name: "João Silva"}

		customerRepo := &mockCustomerRepository{
			FindByIDFunc: func(id string) (*domain.Customer, error) {
				return expected, nil
			},
		}

		customerService := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		req := httptest.NewRequest(http.MethodGet, "/customers/CUST-001", nil)
		req = withURLParam(req, "id", "CUST-001")
		w := httptest.NewRecorder()

		controller.FindCustomerByID(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", w.Result().StatusCode)
		}
	})

	t.Run("deve retornar 404 quando cliente não existe", func(t *testing.T) {
		customerRepo := &mockCustomerRepository{
			FindByIDFunc: func(id string) (*domain.Customer, error) {
				return nil, domain.ErrCustomerNotFound
			},
		}

		customerService := service.NewCustomerService(customerRepo, &mockCustomerIDGenerator{})
		controller := controllers.NewCustomerController(customerService)

		req := httptest.NewRequest(http.MethodGet, "/customers/CUST-999", nil)
		req = withURLParam(req, "id", "CUST-999")
		w := httptest.NewRecorder()

		controller.FindCustomerByID(w, req)

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", w.Result().StatusCode)
		}
	})
}
