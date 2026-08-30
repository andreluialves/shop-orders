package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
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

func TestHandleError_DoesNotLeakInternalDetails(t *testing.T) {
	t.Run("erro desconhecido deve retornar mensagem genérica, sem detalhes internos", func(t *testing.T) {
		w := httptest.NewRecorder()

		sensitiveErr := errors.New("connection refused: postgres://user:senha123@10.0.0.5:5432/shop_orders")

		handleError(w, sensitiveErr)

		resp := w.Result()
		body := w.Body.String()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", resp.StatusCode)
		}

		if strings.Contains(body, "senha123") || strings.Contains(body, "10.0.0.5") || strings.Contains(body, "postgres://") {
			t.Errorf("resposta não deveria conter detalhes internos do erro, recebeu: %v", body)
		}

		if !strings.Contains(body, "internal server error") {
			t.Errorf("esperava mensagem genérica 'internal server error', recebeu: %v", body)
		}
	})

	t.Run("erro de domínio conhecido pode expor sua própria mensagem", func(t *testing.T) {
		w := httptest.NewRecorder()

		handleError(w, domain.ErrCustomerNotFound)

		body := w.Body.String()

		if !strings.Contains(body, domain.ErrCustomerNotFound.Error()) {
			t.Errorf("esperava mensagem de domínio na resposta, recebeu: %v", body)
		}
	})
}
