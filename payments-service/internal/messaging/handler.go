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

func (h *PaymentRequestHandler) Start(ctx context.Context) error {
	msgs, err := h.rabbit.Consume("payment.requested")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event PaymentRequested

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				h.logger.Error("evento payment.requested inválido", "error", err)
				continue
			}

			payment, err := h.paymentService.ProcessPayment(ctx, event.OrderID, event.Amount)
			if err != nil {
				h.logger.Error("falha ao processar pagamento", "order_id", event.OrderID, "error", err)
				continue
			}

			result := PaymentProcessed{
				OrderID:   payment.OrderID,
				PaymentID: payment.ID,
				Status:    string(payment.Status()),
			}

			body, _ := json.Marshal(result)

			if err := h.rabbit.Publish("payment.processed", body); err != nil {
				h.logger.Error("falha ao publicar payment.processed", "order_id", event.OrderID, "error", err)
			}
		}
	}()

	return nil
}
