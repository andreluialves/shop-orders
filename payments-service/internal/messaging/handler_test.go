package messaging_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
	"github.com/andreluialves/shop-orders/payments-service/internal/messaging"
	"github.com/andreluialves/shop-orders/payments-service/internal/service"
)

func TestPaymentRequestHandler_HandleMessage(t *testing.T) {
	t.Run("deve processar solicitação e retornar evento de aprovação", func(t *testing.T) {
		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error { return nil },
		}
		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "PAY-001", nil
			},
		}

		paymentService := service.NewPaymentService(repo, idGen)
		handler := messaging.NewPaymentRequestHandler(nil, paymentService, &mockLogger{})

		event := messaging.PaymentRequested{OrderID: "PED-001", Amount: 500.0}
		body, _ := json.Marshal(event)

		resultBody, err := handler.HandleMessage(context.Background(), body)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		var result messaging.PaymentProcessed
		if err := json.Unmarshal(resultBody, &result); err != nil {
			t.Fatalf("erro ao decodificar resultado: %v", err)
		}

		if result.OrderID != "PED-001" {
			t.Errorf("esperava order_id PED-001, recebeu %v", result.OrderID)
		}

		if result.Status != "APPROVED" {
			t.Errorf("esperava status APPROVED, recebeu %v", result.Status)
		}
	})

	t.Run("deve retornar evento de recusa para valor acima do limite", func(t *testing.T) {
		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error { return nil },
		}
		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "PAY-002", nil
			},
		}

		paymentService := service.NewPaymentService(repo, idGen)
		handler := messaging.NewPaymentRequestHandler(nil, paymentService, &mockLogger{})

		event := messaging.PaymentRequested{OrderID: "PED-002", Amount: 20000.0}
		body, _ := json.Marshal(event)

		resultBody, err := handler.HandleMessage(context.Background(), body)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		var result messaging.PaymentProcessed
		json.Unmarshal(resultBody, &result)

		if result.Status != "DECLINED" {
			t.Errorf("esperava status DECLINED, recebeu %v", result.Status)
		}
	})

	t.Run("deve retornar erro para JSON inválido", func(t *testing.T) {
		paymentService := service.NewPaymentService(&mockPaymentRepository{}, &mockPaymentIDGenerator{})
		handler := messaging.NewPaymentRequestHandler(nil, paymentService, &mockLogger{})

		_, err := handler.HandleMessage(context.Background(), []byte("{json invalido"))

		if err == nil {
			t.Error("esperava erro ao decodificar JSON inválido")
		}
	})

	t.Run("deve propagar erro quando processamento do pagamento falha", func(t *testing.T) {
		paymentService := service.NewPaymentService(&mockPaymentRepository{}, &mockPaymentIDGenerator{})
		handler := messaging.NewPaymentRequestHandler(nil, paymentService, &mockLogger{})

		event := messaging.PaymentRequested{OrderID: "", Amount: 500.0} // orderID vazio → erro de domínio
		body, _ := json.Marshal(event)

		_, err := handler.HandleMessage(context.Background(), body)

		if err == nil {
			t.Error("esperava erro propagado do ProcessPayment")
		}
	})
}
