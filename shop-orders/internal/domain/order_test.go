package domain_test

import (
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
)

func TestNewOrder(t *testing.T) {
	validProduct := &domain.Product{
		ID:       "P001",
		Name:     "Notebook",
		Price:    3500,
		Quantity: 10,
	}

	tests := []struct {
		name       string
		customerID string
		items      []*domain.OrderItem
		wantErr    error
	}{
		{
			name:       "deve criar pedido válido",
			customerID: "CUST-001",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: nil,
		},
		{
			name:       "não deve criar pedido sem customerID",
			customerID: "",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: domain.ErrInvalidCustomer,
		},
		{
			name:       "não deve criar pedido com customerID só com espaços",
			customerID: "   ",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 2, 3500),
			},
			wantErr: domain.ErrInvalidCustomer,
		},
		{
			name:       "não deve criar pedido sem itens",
			customerID: "CUST-001",
			items:      []*domain.OrderItem{},
			wantErr:    domain.ErrEmptyOrder,
		},
		{
			name:       "não deve criar pedido com item sem produto",
			customerID: "CUST-001",
			items: []*domain.OrderItem{
				domain.NewOrderItem(nil, 2, 3500),
			},
			wantErr: domain.ErrProductNotFound,
		},
		{
			name:       "não deve criar pedido com quantidade inválida",
			customerID: "CUST-001",
			items: []*domain.OrderItem{
				domain.NewOrderItem(validProduct, 0, 3500),
			},
			wantErr: domain.ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := domain.NewOrder("ORD-001", tt.customerID, tt.items)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				if order != nil {
					t.Errorf("esperava order nil em caso de erro, recebeu %v", order)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if order.CustomerID != tt.customerID {
				t.Errorf("esperava customerID %v, recebeu %v", tt.customerID, order.CustomerID)
			}

			if order.Status() != domain.OrderStatusPending {
				t.Errorf("esperava status PENDING, recebeu %v", order.Status())
			}
		})
	}
}

func TestOrder_Pay(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus domain.OrderStatus
		wantErr       error
	}{
		{
			name:          "deve pagar pedido pendente",
			initialStatus: domain.OrderStatusPending,
			wantErr:       nil,
		},
		{
			name:          "não deve pagar pedido já pago",
			initialStatus: domain.OrderStatusPaid,
			wantErr:       domain.ErrChangeStatusInvalid,
		},
		{
			name:          "não deve pagar pedido cancelado",
			initialStatus: domain.OrderStatusCanceled,
			wantErr:       domain.ErrChangeStatusInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := domain.RestoreOrder("PED-001", "CUST-001", tt.initialStatus)

			err := order.Pay()

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if order.Status() != domain.OrderStatusPaid {
				t.Errorf("esperava status PAID, recebeu %v", order.Status())
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := domain.RestoreOrder("PED-001", "CUST-001", tt.initialStatus)

			err := order.Cancel()

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
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

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := domain.NewOrder("ORD-001", "CUST-001", tt.items)

			if err != nil {
				order = domain.RestoreOrder("ORD-001", "CUST-001", domain.OrderStatusPending)
				for _, item := range tt.items {
					order.AddItem(item)
				}
			}

			got := order.TotalSum()

			if got != tt.want {
				t.Errorf("esperava total %v, recebeu %v", tt.want, got)
			}
		})
	}
}

func TestOrderItem_Validate(t *testing.T) {
	validProduct := &domain.Product{ID: "P001", Name: "Notebook", Price: 3500, Quantity: 10}

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()

			if err != tt.wantErr {
				t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
			}
		})
	}
}
