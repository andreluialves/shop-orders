package domain_test

import (
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
)

func TestNewProduct(t *testing.T) {
	testCases := []struct {
		name     string
		product  string
		price    float64
		quantity int
		wantErr  error
	}{
		{
			name:     "deve criar produto válido",
			product:  "Notebook",
			price:    3500,
			quantity: 10,
			wantErr:  nil,
		},
		{
			name:     "não deve criar produto sem nome",
			product:  "",
			price:    3500,
			quantity: 10,
			wantErr:  domain.ErrProductNameInvalid,
		},
		{
			name:     "não deve criar produto com preço zero",
			product:  "Notebook",
			price:    0,
			quantity: 10,
			wantErr:  domain.ErrProductPriceInvalid,
		},
		{
			name:     "não deve criar produto com preço negativo",
			product:  "Notebook",
			price:    -100,
			quantity: 10,
			wantErr:  domain.ErrProductPriceInvalid,
		},
		{
			name:     "não deve criar produto com quantidade negativa",
			product:  "Notebook",
			price:    3500,
			quantity: -1,
			wantErr:  domain.ErrInvalidQuantity,
		},
		{
			name:     "deve criar produto com quantidade zero",
			product:  "Notebook",
			price:    3500,
			quantity: 0,
			wantErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product, err := domain.NewProduct(tc.product, tc.price, tc.quantity)

			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
				}
				if product != nil {
					t.Errorf("esperava product nil em caso de erro, recebeu %v", product)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if product.Name != tc.product {
				t.Errorf("esperava name %v, recebeu %v", tc.product, product.Name)
			}

			if product.Quantity != tc.quantity {
				t.Errorf("esperava quantity %v, recebeu %v", tc.quantity, product.Quantity)
			}
		})
	}
}

func TestProduct_ReduceQuantity(t *testing.T) {
	testCases := []struct {
		name         string
		stock        int
		reduceBy     int
		wantErr      error
		wantQuantity int
	}{
		{
			name:         "deve reduzir quando há estoque suficiente",
			stock:        10,
			reduceBy:     4,
			wantErr:      nil,
			wantQuantity: 6,
		},
		{
			name:         "deve reduzir quando reduz todo o estoque",
			stock:        10,
			reduceBy:     10,
			wantErr:      nil,
			wantQuantity: 0,
		},
		{
			name:         "não deve reduzir quando quantidade excede estoque",
			stock:        5,
			reduceBy:     10,
			wantErr:      domain.ErrInsufficientQuantity,
			wantQuantity: 5,
		},
		{
			name:         "não deve reduzir quando quantidade é zero",
			stock:        10,
			reduceBy:     0,
			wantErr:      domain.ErrInvalidQuantity,
			wantQuantity: 10,
		},
		{
			name:         "não deve reduzir quando quantidade é negativa",
			stock:        10,
			reduceBy:     -1,
			wantErr:      domain.ErrInvalidQuantity,
			wantQuantity: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product := domain.RestoreProduct("P001", "Notebook", 3500, tc.stock)

			err := product.ReduceQuantity(tc.reduceBy)

			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if product.Quantity != tc.wantQuantity {
				t.Errorf("esperava quantity %v, recebeu %v", tc.wantQuantity, product.Quantity)
			}
		})
	}
}

func TestProduct_RestoreQuantity(t *testing.T) {
	testCases := []struct {
		name         string
		stock        int
		restoreBy    int
		wantErr      error
		wantQuantity int
	}{
		{
			name:         "deve restaurar quantidade normalmente",
			stock:        10,
			restoreBy:    5,
			wantErr:      nil,
			wantQuantity: 15,
		},
		{
			name:         "não deve restaurar quantidade zero",
			stock:        10,
			restoreBy:    0,
			wantErr:      domain.ErrInvalidQuantity,
			wantQuantity: 10,
		},
		{
			name:         "não deve restaurar quantidade negativa",
			stock:        10,
			restoreBy:    -5,
			wantErr:      domain.ErrInvalidQuantity,
			wantQuantity: 10,
		},
		{
			name:         "restaurar quantidade maior que o estoque atual retorna erro",
			stock:        2,
			restoreBy:    5,
			wantErr:      domain.ErrInsufficientQuantity,
			wantQuantity: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product := domain.RestoreProduct("P001", "Notebook", 3500, tc.stock)

			err := product.RestoreQuantity(tc.restoreBy)

			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tc.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if product.Quantity != tc.wantQuantity {
				t.Errorf("esperava quantity %v, recebeu %v", tc.wantQuantity, product.Quantity)
			}
		})
	}
}
