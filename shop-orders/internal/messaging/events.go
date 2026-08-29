package messaging

// PaymentRequested é publicado pelo serviço de orders quando um pagamento
// é solicitado. Consumido pelo serviço de payments.
type PaymentRequested struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

// PaymentProcessed é publicado pelo serviço de payments quando um
// pagamento termina de ser processado (aprovado ou recusado).
type PaymentProcessed struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"` // "APPROVED" ou "DECLINED"
}
