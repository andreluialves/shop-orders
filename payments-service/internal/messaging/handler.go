package messaging

import (
	"context"
	"encoding/json"

	"github.com/andreluialves/shop-orders/payments-service/internal/logger"
	"github.com/andreluialves/shop-orders/payments-service/internal/service"
)

type PaymentRequestHandler struct {
	rabbit         *RabbitMQ
	paymentService *service.PaymentService
	logger         logger.Logger
}

func NewPaymentRequestHandler(rabbit *RabbitMQ, paymentService *service.PaymentService, logger logger.Logger) *PaymentRequestHandler {
	return &PaymentRequestHandler{rabbit: rabbit, paymentService: paymentService, logger: logger}
}

// HandleMessage contém a lógica de negócio pura: decodifica o evento,
// processa o pagamento, monta o evento de resposta. Não depende do
// RabbitMQ, o que permite testar sem um broker real.
func (h *PaymentRequestHandler) HandleMessage(ctx context.Context, body []byte) ([]byte, error) {
	var event PaymentRequested

	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}

	payment, err := h.paymentService.ProcessPayment(ctx, event.OrderID, event.Amount)
	if err != nil {
		return nil, err
	}

	result := PaymentProcessed{
		OrderID:   payment.OrderID,
		PaymentID: payment.ID,
		Status:    string(payment.Status()),
	}

	return json.Marshal(result)
}

// Start cuida só da mecânica do RabbitMQ: consumir, delegar pra
// HandleMessage, publicar o resultado.
func (h *PaymentRequestHandler) Start(ctx context.Context) error {
	msgs, err := h.rabbit.Consume("payment.requested")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			resultBody, err := h.HandleMessage(ctx, msg.Body)
			if err != nil {
				h.logger.Error("falha ao processar payment.requested", "error", err)
				continue
			}

			if err := h.rabbit.Publish("payment.processed", resultBody); err != nil {
				h.logger.Error("falha ao publicar payment.processed", "error", err)
			}
		}
	}()

	return nil
}
