package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreluialves/shop-orders/internal/controllers"
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/dto"
	"github.com/andreluialves/shop-orders/internal/service"
)

func TestOrderController_CreateOrder(t *testing.T) {
	t.Run("deve criar pedido com sucesso e retornar 201", func(t *testing.T) {
		product := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return product, nil
			},
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		orderRepo := &mockOrderRepository{
			SaveFunc: func(o *domain.Order) error {
				return nil
			},
		}

		orderService := service.NewOrderService(productRepo, orderRepo)
		controller := controllers.NewOrderController(orderService)

		body := `{"customer":"João Silva","items":[{"id":"P001","quantity":2}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("esperava status 201, recebeu %d", resp.StatusCode)
		}

		var got dto.OrderResponse
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.Customer != "João Silva" {
			t.Errorf("esperava customer João Silva, recebeu %v", got.Customer)
		}

		if got.Status != string(domain.OrderStatusPending) {
			t.Errorf("esperava status PENDING, recebeu %v", got.Status)
		}

		if len(got.Items) != 1 {
			t.Errorf("esperava 1 item, recebeu %d", len(got.Items))
		}
	})

	t.Run("deve retornar 400 quando body é JSON inválido", func(t *testing.T) {
		orderService := service.NewOrderService(&mockProductRepository{}, &mockOrderRepository{})
		controller := controllers.NewOrderController(orderService)

		body := `{"customer": "João",` // malformado
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 404 quando produto não existe", func(t *testing.T) {
		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}

		orderService := service.NewOrderService(productRepo, &mockOrderRepository{})
		controller := controllers.NewOrderController(orderService)

		body := `{"customer":"João Silva","items":[{"id":"P999","quantity":2}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 409 quando quantidade solicitada excede estoque", func(t *testing.T) {
		product := domain.RestoreProduct("P001", "Notebook", 3500, 1)

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return product, nil
			},
		}

		orderService := service.NewOrderService(productRepo, &mockOrderRepository{})
		controller := controllers.NewOrderController(orderService)

		body := `{"customer":"João Silva","items":[{"id":"P001","quantity":5}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("esperava status 409, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 400 quando customer é vazio", func(t *testing.T) {
		product := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return product, nil
			},
		}

		orderService := service.NewOrderService(productRepo, &mockOrderRepository{})
		controller := controllers.NewOrderController(orderService)

		body := `{"customer":"","items":[{"id":"P001","quantity":2}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
		}
	})
}

func TestOrderController_FindAllOrders(t *testing.T) {
	t.Run("deve retornar lista paginada com 200", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

		orderRepo := &mockOrderRepository{
			ListFunc: func(limit, offset int) ([]*domain.Order, int, error) {
				return []*domain.Order{order}, 1, nil
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodGet, "/orders?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		controller.FindAllOrders(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var body struct {
			Data  []dto.OrderResponse `json:"data"`
			Total int                 `json:"total"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if body.Total != 1 {
			t.Errorf("esperava total 1, recebeu %d", body.Total)
		}

		if len(body.Data) != 1 {
			t.Errorf("esperava 1 pedido, recebeu %d", len(body.Data))
		}
	})

	t.Run("deve retornar 500 quando List falha", func(t *testing.T) {
		orderRepo := &mockOrderRepository{
			ListFunc: func(limit, offset int) ([]*domain.Order, int, error) {
				return nil, 0, errors.New("erro de conexão com banco")
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		controller.FindAllOrders(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", resp.StatusCode)
		}
	})
}

func TestOrderController_FindOrderByID(t *testing.T) {
	t.Run("deve retornar pedido com 200 quando encontrado", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodGet, "/orders/PED-001", nil)
		req = withURLParam(req, "id", "PED-001")
		w := httptest.NewRecorder()

		controller.FindOrderByID(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 404 quando pedido não existe", func(t *testing.T) {
		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return nil, domain.ErrOrderNotFound
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodGet, "/orders/PED-999", nil)
		req = withURLParam(req, "id", "PED-999")
		w := httptest.NewRecorder()

		controller.FindOrderByID(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", resp.StatusCode)
		}
	})
}

func TestOrderController_PayOrder(t *testing.T) {
	t.Run("deve pagar pedido pendente e retornar 200", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				return nil
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-001/pay", nil)
		req = withURLParam(req, "id", "PED-001")
		w := httptest.NewRecorder()

		controller.PayOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var got dto.OrderResponse
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.Status != string(domain.OrderStatusPaid) {
			t.Errorf("esperava status PAID na resposta, recebeu %v", got.Status)
		}
	})

	t.Run("deve retornar 409 quando pedido já está pago", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPaid)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-001/pay", nil)
		req = withURLParam(req, "id", "PED-001")
		w := httptest.NewRecorder()

		controller.PayOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("esperava status 409, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 404 quando pedido não existe", func(t *testing.T) {
		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return nil, domain.ErrOrderNotFound
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-999/pay", nil)
		req = withURLParam(req, "id", "PED-999")
		w := httptest.NewRecorder()

		controller.PayOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", resp.StatusCode)
		}
	})
}

func TestOrderController_CancelOrder(t *testing.T) {
	t.Run("deve cancelar pedido pendente e restaurar estoque, retornando 200", func(t *testing.T) {
		product := domain.RestoreProduct("P001", "Notebook", 3500, 5)
		item := domain.NewOrderItem(product, 2, 3500)

		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)
		order.AddItem(item)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				return nil
			},
		}

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return domain.RestoreProduct("P001", "Notebook", 3500, 5), nil
			},
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		orderService := service.NewOrderService(productRepo, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-001/cancel", nil)
		req = withURLParam(req, "id", "PED-001")
		w := httptest.NewRecorder()

		controller.CancelOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var got dto.OrderResponse
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.Status != string(domain.OrderStatusCanceled) {
			t.Errorf("esperava status CANCELED na resposta, recebeu %v", got.Status)
		}
	})

	t.Run("deve retornar 409 quando pedido já está cancelado", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusCanceled)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-001/cancel", nil)
		req = withURLParam(req, "id", "PED-001")
		w := httptest.NewRecorder()

		controller.CancelOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("esperava status 409, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 404 quando pedido não existe", func(t *testing.T) {
		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return nil, domain.ErrOrderNotFound
			},
		}

		orderService := service.NewOrderService(&mockProductRepository{}, orderRepo)
		controller := controllers.NewOrderController(orderService)

		req := httptest.NewRequest(http.MethodPost, "/orders/PED-999/cancel", nil)
		req = withURLParam(req, "id", "PED-999")
		w := httptest.NewRecorder()

		controller.CancelOrder(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", resp.StatusCode)
		}
	})
}
