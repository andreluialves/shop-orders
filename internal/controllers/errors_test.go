package controllers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
)

func TestHandleError(t *testing.T) {
	testcases := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"produto não encontrado", domain.ErrProductNotFound, 404},
		{"pedido não encontrado", domain.ErrOrderNotFound, 404},
		{"quantidade insuficiente", domain.ErrInsufficientQuantity, 409},
		{"mudança de status inválida", domain.ErrChangeStatusInvalid, 409},
		{"pedido já pago", domain.ErrOrderAlreadyPaid, 409},
		{"pedido já cancelado", domain.ErrOrderAlreadyCanceled, 409},
		{"cliente inválido", domain.ErrInvalidCustomer, 400},
		{"pedido vazio", domain.ErrEmptyOrder, 400},
		{"nome do produto inválido", domain.ErrProductNameInvalid, 400},
		{"preço do produto inválido", domain.ErrProductPriceInvalid, 400},
		{"erro desconhecido", errors.New("algo inesperado"), 500},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handleError(w, tc.err)

			if w.Result().StatusCode != tc.wantStatus {
				t.Errorf("esperava status %d, recebeu %d", tc.wantStatus, w.Result().StatusCode)
			}
		})
	}
}
