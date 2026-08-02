package service_test

import (
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/service"
)

func TestOrderService_PayOrder(t *testing.T) {
	t.Run("deve pagar pedido pendente com sucesso", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

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

		s := service.NewOrderService(&mockProductRepository{}, orderRepo)

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

		s := service.NewOrderService(&mockProductRepository{}, orderRepo)

		_, err := s.PayOrder("PED-999")

		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Errorf("esperava ErrOrderNotFound, recebeu %v", err)
		}
	})

	t.Run("não deve pagar pedido já pago, e Update não deve ser chamado", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPaid)

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

		s := service.NewOrderService(&mockProductRepository{}, orderRepo)

		_, err := s.PayOrder("PED-001")

		if !errors.Is(err, domain.ErrChangeStatusInvalid) {
			t.Errorf("esperava ErrChangeStatusInvalid, recebeu %v", err)
		}

		if updateCalled {
			t.Error("Update não deveria ter sido chamado quando Pay() falha")
		}
	})

	t.Run("deve propagar erro quando Update falha", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

		orderRepo := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return order, nil
			},
			UpdateFunc: func(o *domain.Order) error {
				return errors.New("erro de conexão com banco")
			},
		}

		s := service.NewOrderService(&mockProductRepository{}, orderRepo)

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
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)
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

		s := service.NewOrderService(productRepo, orderRepo)

		result, err := s.CancelOrder("PED-001")

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
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusCanceled)
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

		s := service.NewOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder("PED-001")

		if !errors.Is(err, domain.ErrChangeStatusInvalid) {
			t.Errorf("esperava ErrChangeStatusInvalid, recebeu %v", err)
		}

		if productSaveCalled {
			t.Error("estoque não deveria ser tocado quando Cancel() falha antes do loop")
		}
	})

	t.Run("deve retornar erro quando produto não é encontrado no loop", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)
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

		s := service.NewOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder("PED-001")

		if err == nil {
			t.Error("esperava erro quando produto não é encontrado, recebeu nil")
		}
	})

	t.Run("não deve chamar Update do pedido se falhar ao salvar produto", func(t *testing.T) {
		order := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)
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

		s := service.NewOrderService(productRepo, orderRepo)

		_, err := s.CancelOrder("PED-001")

		if err == nil {
			t.Error("esperava erro propagado do Save do produto")
		}

		if orderUpdateCalled {
			t.Error("Update do pedido não deveria ser chamado se Save do produto falhou antes")
		}
	})
}
