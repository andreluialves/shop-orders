package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/shop-orders/internal/service"
)

func TestOrderService_PayOrder(t *testing.T) {
	t.Run("deve pagar pedido pendente com sucesso", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPending)

		var updateCalled bool

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				updateCalled = true
				return nil
			},
		}

		s := newTestOrderService(&mockProductRepository{}, orderRepo)

		result, err := s.PayOrder("PED-001")

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if result.Status() != domain.OrderStatusPaid {
			t.Errorf("esperava status PAID, recebeu %v", result.Status())
		}

		if !updateCalled {
			t.Error("esperava que Update fosse chamado, mas não foi")
		}
	})

	t.Run("deve retornar erro quando pedido não existe", func(t *testing.T) {
		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return nil, domain.ErrOrderNotFound
			},
		}

		s := newTestOrderService(&mockProductRepository{}, orderRepo)

		_, err := s.PayOrder("PED-999")

		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Errorf("esperava ErrOrderNotFound, recebeu %v", err)
		}
	})

	t.Run("não deve pagar pedido já pago, e Update não deve ser chamado", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPaid)

		var updateCalled bool

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				updateCalled = true
				return nil
			},
		}

		s := newTestOrderService(&mockProductRepository{}, orderRepo)

		_, err := s.PayOrder("PED-001")

		if !errors.Is(err, domain.ErrChangeStatusInvalid) {
			t.Errorf("esperava ErrChangeStatusInvalid, recebeu %v", err)
		}

		if updateCalled {
			t.Error("Update não deveria ter sido chamado quando Pay() falha")
		}
	})

	t.Run("deve propagar erro quando Update falha", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPending)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				return errors.New("erro de conexão com banco")
			},
		}

		s := newTestOrderService(&mockProductRepository{}, orderRepo)

		_, err := s.PayOrder("PED-001")

		if err == nil {
			t.Error("esperava erro propagado do Update, recebeu nil")
		}
	})
}

func TestOrderService_CancelOrder(t *testing.T) {
	product := domain.RestoreProduct("P001", "Notebook", 3500, 5)
	item := domain.NewOrderItem(product, 2, 3500)

	t.Run("deve cancelar pedido e restaurar estoque", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPending)
		order.AddItem(item)

		var savedQuantity int
		var updateCalled bool

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				updateCalled = true
				return nil
			},
		}

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return domain.RestoreProduct("P001", "Notebook", 3500, 5), nil
			},
			SaveFunc: func(p *domain.Product) error {
				savedQuantity = p.Quantity
				return nil
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		result, err := s.CancelOrder(context.Background(), "PED-001")

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if result.Status() != domain.OrderStatusCanceled {
			t.Errorf("esperava status CANCELED, recebeu %v", result.Status())
		}

		if savedQuantity != 7 {
			t.Errorf("esperava estoque restaurado para 7 (5+2), recebeu %v", savedQuantity)
		}

		if !updateCalled {
			t.Error("esperava que Update do pedido fosse chamado")
		}
	})

	t.Run("não deve cancelar pedido já cancelado", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusCanceled)
		order.AddItem(item)

		var productSaveCalled bool

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
		}

		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				productSaveCalled = true
				return nil
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder(context.Background(), "PED-001")

		if !errors.Is(err, domain.ErrChangeStatusInvalid) {
			t.Errorf("esperava ErrChangeStatusInvalid, recebeu %v", err)
		}

		if productSaveCalled {
			t.Error("estoque não deveria ser tocado quando Cancel() falha antes do loop")
		}
	})

	t.Run("deve retornar erro quando produto não é encontrado no loop", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPending)
		order.AddItem(item)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
		}

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, errors.New("produto não encontrado")
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder(context.Background(), "PED-001")

		if err == nil {
			t.Error("esperava erro quando produto não é encontrado, recebeu nil")
		}
	})

	t.Run("não deve chamar Update do pedido se falhar ao salvar produto", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "CUST-001", domain.OrderStatusPending)
		order.AddItem(item)

		var orderUpdateCalled bool

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				orderUpdateCalled = true
				return nil
			},
		}

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return domain.RestoreProduct("P001", "Notebook", 3500, 5), nil
			},
			SaveFunc: func(p *domain.Product) error {
				return errors.New("falha ao salvar produto")
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder(context.Background(), "PED-001")

		if err == nil {
			t.Error("esperava erro propagado do Save do produto")
		}

		if orderUpdateCalled {
			t.Error("Update do pedido não deveria ser chamado se Save do produto falhou antes")
		}
	})
}

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("deve criar pedido quando cliente e produtos existem", func(t *testing.T) {
		product := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		var savedOrder *domain.Order
		var savedProductQuantity int

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return product, nil
			},
			SaveFunc: func(p *domain.Product) error {
				savedProductQuantity = p.Quantity
				return nil
			},
		}

		orderRepo := &mockOrderRepository{
			SaveFunc: func(o *domain.Order) error {
				savedOrder = o
				return nil
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		items := []service.CreateOrderItem{{ID: "P001", Quantity: 2}}

		order, err := s.CreateOrder(context.Background(), "CUST-001", items)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if order.CustomerID != "CUST-001" {
			t.Errorf("esperava CustomerID CUST-001, recebeu %v", order.CustomerID)
		}

		if savedProductQuantity != 8 {
			t.Errorf("esperava estoque reduzido para 8 (10-2), recebeu %v", savedProductQuantity)
		}

		if savedOrder == nil {
			t.Error("esperava que Order.Save fosse chamado")
		}
	})

	t.Run("não deve criar pedido quando cliente não existe", func(t *testing.T) {
		var orderSaveCalled bool
		var productFindCalled bool

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				productFindCalled = true
				return domain.RestoreProduct("P001", "Notebook", 3500, 10), nil
			},
		}

		orderRepo := &mockOrderRepository{
			SaveFunc: func(o *domain.Order) error {
				orderSaveCalled = true
				return nil
			},
		}

		customerRepo := &mockCustomerRepository{
			FindByIDFunc: func(id string) (*domain.Customer, error) {
				return nil, domain.ErrCustomerNotFound
			},
		}

		s := newTestOrderServiceWithCustomerRepo(productRepo, orderRepo, customerRepo)

		items := []service.CreateOrderItem{{ID: "P001", Quantity: 2}}

		_, err := s.CreateOrder(context.Background(), "CUST-999", items)

		if !errors.Is(err, domain.ErrCustomerNotFound) {
			t.Errorf("esperava ErrCustomerNotFound, recebeu %v", err)
		}

		if productFindCalled {
			t.Error("não deveria buscar produtos se o cliente não existe")
		}

		if orderSaveCalled {
			t.Error("Order.Save não deveria ser chamado quando cliente não existe")
		}
	})

	t.Run("não deve salvar pedido se produto não é encontrado", func(t *testing.T) {
		var orderSaveCalled bool

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}

		orderRepo := &mockOrderRepository{
			SaveFunc: func(o *domain.Order) error {
				orderSaveCalled = true
				return nil
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		items := []service.CreateOrderItem{{ID: "P999", Quantity: 2}}

		_, err := s.CreateOrder(context.Background(), "CUST-001", items)

		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Errorf("esperava ErrProductNotFound, recebeu %v", err)
		}

		if orderSaveCalled {
			t.Error("Order.Save não deveria ser chamado quando produto não é encontrado")
		}
	})

	t.Run("não deve salvar pedido se estoque de algum item é insuficiente", func(t *testing.T) {
		productA := domain.RestoreProduct("P001", "Notebook", 3500, 10)
		productB := domain.RestoreProduct("P002", "Mouse", 150, 1)

		var orderSaveCalled bool

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				if id == "P001" {
					return productA, nil
				}
				return productB, nil
			},
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		orderRepo := &mockOrderRepository{
			SaveFunc: func(o *domain.Order) error {
				orderSaveCalled = true
				return nil
			},
		}

		s := newTestOrderService(productRepo, orderRepo)

		items := []service.CreateOrderItem{
			{ID: "P001", Quantity: 2},
			{ID: "P002", Quantity: 5},
		}

		_, err := s.CreateOrder(context.Background(), "CUST-001", items)

		if !errors.Is(err, domain.ErrInsufficientQuantity) {
			t.Errorf("esperava ErrInsufficientQuantity, recebeu %v", err)
		}

		if orderSaveCalled {
			t.Error("Order.Save não deveria ser chamado quando algum item tem estoque insuficiente")
		}
	})
}
