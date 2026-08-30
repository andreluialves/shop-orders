package events

// EventPublisher é o contrato que o service depende para publicar eventos,
// sem conhecer os detalhes de RabbitMQ (ou qualquer outro broker).
type EventPublisher interface {
	Publish(queue string, body []byte) error
}

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
