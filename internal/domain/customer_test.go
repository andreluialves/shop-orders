package domain_test

import (
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
)

func TestNewCustomer(t *testing.T) {
	tests := []struct {
		name         string
		customerName string
		email        string
		address      string
		phone        string
		wantErr      error
	}{
		{
			name:         "deve criar cliente válido",
			customerName: "João Silva",
			email:        "joao@example.com",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      nil,
		},
		{
			name:         "não deve criar cliente sem nome",
			customerName: "",
			email:        "joao@example.com",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      domain.ErrCustomerNameRequired,
		},
		{
			name:         "não deve criar cliente sem email",
			customerName: "João Silva",
			email:        "",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      domain.ErrCustomerEmailRequired,
		},
		{
			name:         "não deve criar cliente com email sem @",
			customerName: "João Silva",
			email:        "joaoexample.com",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      domain.ErrCustomerEmailInvalid,
		},
		{
			name:         "não deve criar cliente com email sem domínio",
			customerName: "João Silva",
			email:        "joao@",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      domain.ErrCustomerEmailInvalid,
		},
		{
			name:         "deve aceitar email com subdomínio e sinal de mais",
			customerName: "João Silva",
			email:        "joao+pedidos@sub.dominio.com.br",
			address:      "Rua A, 123",
			phone:        "11999999999",
			wantErr:      nil,
		},
		{
			name:         "não deve criar cliente sem endereço",
			customerName: "João Silva",
			email:        "joao@example.com",
			address:      "",
			phone:        "11999999999",
			wantErr:      domain.ErrCustomerAddressRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := domain.NewCustomer(tt.customerName, tt.email, tt.address, tt.phone)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				if customer != nil {
					t.Errorf("esperava customer nil em caso de erro, recebeu %v", customer)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if customer.Name != tt.customerName {
				t.Errorf("esperava name %v, recebeu %v", tt.customerName, customer.Name)
			}

			if customer.Email != tt.email {
				t.Errorf("esperava email %v, recebeu %v", tt.email, customer.Email)
			}
		})
	}
}
