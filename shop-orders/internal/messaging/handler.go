package messaging

import (
	"context"
	"encoding/json"

	"github.com/andreluialves/shop-orders/shop-orders/internal/events"
	"github.com/andreluialves/shop-orders/shop-orders/internal/logger"
	"github.com/andreluialves/shop-orders/shop-orders/internal/service"
)

type PaymentProcessedHandler struct {
	rabbit       *RabbitMQ
	orderService *service.OrderService
	logger       logger.Logger
}

func NewPaymentProcessedHandler(rabbit *RabbitMQ, orderService *service.OrderService, logger logger.Logger) *PaymentProcessedHandler {
	return &PaymentProcessedHandler{rabbit: rabbit, orderService: orderService, logger: logger}
}

func (h *PaymentProcessedHandler) Start(ctx context.Context) error {
	msgs, err := h.rabbit.Consume("payment.processed")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event events.PaymentProcessed

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				h.logger.Error("evento payment.processed inválido", "error", err)
				continue
			}

			if err := h.orderService.HandlePaymentResult(ctx, event.OrderID, event.Status); err != nil {
				h.logger.Error("falha ao tratar resultado de pagamento", "order_id", event.OrderID, "error", err)
			}
		}
	}()

	return nil
}
