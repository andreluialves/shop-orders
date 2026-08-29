package dto

import "github.com/andreluialves/shop-orders/payments-service/internal/domain"

type PaymentResponse struct {
	ID      string  `json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}

func NewPaymentResponse(payment *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		Amount:  payment.Amount,
		Status:  string(payment.Status()),
	}
}
