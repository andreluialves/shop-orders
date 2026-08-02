package domain_test

import (
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
)

func TestNewOrder(t *testing.T) {
	validProduct := &domain.Product{
		ID:       "P001",
		Name:     "Notebook",
		Price:    3500,
		Quantity: 10,
	}

	testCases := []struct {
		name     string
		customer string
		items    []*domain.OrderItem
		wantErr  error
	}{
		{
			name:     "deve criar pedido válido",
			customer: "João Silva",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: nil,
		},
		{
			name:     "não deve criar pedido sem customer",
			customer: "",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: domain.ErrInvalidCustomer,
		},
		{
			name:     "não deve criar pedido com customer só com espaços",
			customer: "   ",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: domain.ErrInvalidCustomer,
		},
		{
			name:     "não deve criar pedido sem itens",
			customer: "João Silva",
			items:    []*domain.OrderItem{},
			wantErr:  domain.ErrEmptyOrder,
		},
		{
			name:     "não deve criar pedido com item sem produto",
			customer: "João Silva",
			items: []*domain.OrderItem{
				domain.NewOrderItem(nil, 2, 3500),
			},
			wantErr: domain.ErrProductNotFound,
		},
		{
			name:     "não deve criar pedido com quantidade inválida",
			customer: "João Silva",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 0, 3500),
			},
			wantErr: domain.ErrInvalidQuantity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := domain.NewOrder("ORD-001", tc.customer, tc.items)

			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
				}
				if order != nil {
					t.Errorf("esperava order nil em caso de erro, recebeu %v", order)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if order.Status() != domain.OrderStatusPending {
				t.Errorf("esperava status PENDING, recebeu %v", order.Status())
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	testCases := []struct {
		name          string
		initialStatus domain.OrderStatus
		wantErr       error
	}{
		{
			name:          "deve cancelar pedido pendente",
			initialStatus: domain.OrderStatusPending,
			wantErr:       nil,
		},
		{
			name:          "não deve cancelar pedido já pago",
			initialStatus: domain.OrderStatusPaid,
			wantErr:       domain.ErrChangeStatusInvalid,
		},
		{
			name:          "não deve cancelar pedido já cancelado",
			initialStatus: domain.OrderStatusCanceled,
			wantErr:       domain.ErrChangeStatusInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := domain.RestoreOrder("ORD-001", "João Silva", tc.initialStatus)

			err := order.Cancel()

			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if order.Status() != domain.OrderStatusCanceled {
				t.Errorf("esperava status CANCELED, recebeu %v", order.Status())
			}
		})
	}
}

func TestOrder_TotalSum(t *testing.T) {
	product1 := &domain.Product{ID: "P001", Name: "Notebook", Price: 3500, Quantity: 10}
	product2 := &domain.Product{ID: "P002", Name: "Mouse", Price: 150, Quantity: 20}

	testCases := []struct {
		name  string
		items []*domain.OrderItem
		want  float64
	}{
		{
			name: "deve somar corretamente um único item",
			items: []*domain.OrderItem{
				domain.NewOrderItem(product1, 2, 3500),
			},
			want: 7000,
		},
		{
			name: "deve somar corretamente múltiplos itens",
			items: []*domain.OrderItem{
				domain.NewOrderItem(product1, 2, 3500),
				domain.NewOrderItem(product2, 3, 150),
			},
			want: 7450,
		},
		{
			name:  "deve retornar zero sem itens",
			items: []*domain.OrderItem{},
			want:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := domain.NewOrder("ORD-001", "João Silva", tc.items)

			// Para o caso de itens vazios, NewOrder retorna erro (ErrEmptyOrder),
			// então foi usado TotalSum via RestoreOrder + AddItem manual
			if err != nil {
				order = domain.RestoreOrder("ORD-001", "João Silva", domain.OrderStatusPending)
				for _, item := range tc.items {
					order.AddItem(item)
				}
			}

			got := order.TotalSum()

			if got != tc.want {
				t.Errorf("esperava total %v, recebeu %v", tc.want, got)
			}
		})
	}
}

func TestOrderItem_Validate(t *testing.T) {
	validProduct := &domain.Product{ID: "P001", Name: "Notebook", Price: 3500, Quantity: 10}

	testCases := []struct {
		name    string
		item    *domain.OrderItem
		wantErr error
	}{
		{
			name:    "deve validar item correto",
			item:    domain.NewOrderItem(validProduct, 2, 3500),
			wantErr: nil,
		},
		{
			name:    "não deve validar sem produto",
			item:    domain.NewOrderItem(nil, 2, 3500),
			wantErr: domain.ErrProductNotFound,
		},
		{
			name:    "não deve validar com quantidade zero",
			item:    domain.NewOrderItem(validProduct, 0, 3500),
			wantErr: domain.ErrInvalidQuantity,
		},
		{
			name:    "não deve validar com quantidade negativa",
			item:    domain.NewOrderItem(validProduct, -1, 3500),
			wantErr: domain.ErrInvalidQuantity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.item.Validate()

			if err != tc.wantErr {
				t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
			}
		})
	}
}
