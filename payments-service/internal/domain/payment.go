package domain

import "errors"

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusApproved PaymentStatus = "APPROVED"
	PaymentStatusDeclined PaymentStatus = "DECLINED"
)

var (
	ErrPaymentNotFound    = errors.New("pagamento não encontrado")
	ErrInvalidOrderID     = errors.New("order id inválido")
	ErrInvalidAmount      = errors.New("valor de pagamento inválido")
	ErrPaymentAlreadyDone = errors.New("pagamento já foi processado")
)

type Payment struct {
	ID      string
	OrderID string
	Amount  float64
	status  PaymentStatus
}

func NewPayment(id, orderID string, amount float64) (*Payment, error) {
	p := Payment{ID: id, OrderID: orderID, Amount: amount, status: PaymentStatusPending}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p Payment) Validate() error {
	if p.OrderID == "" {
		return ErrInvalidOrderID
	}
	if p.Amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func RestorePayment(id, orderID string, amount float64, status PaymentStatus) *Payment {
	return &Payment{ID: id, OrderID: orderID, Amount: amount, status: status}
}

func (p *Payment) Approve() error {
	if p.status != PaymentStatusPending {
		return ErrPaymentAlreadyDone
	}
	p.status = PaymentStatusApproved
	return nil
}

func (p *Payment) Decline() error {
	if p.status != PaymentStatusPending {
		return ErrPaymentAlreadyDone
	}
	p.status = PaymentStatusDeclined
	return nil
}

func (p Payment) Status() PaymentStatus {
	return p.status
}
