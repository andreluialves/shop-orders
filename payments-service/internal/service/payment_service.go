package service

import (
	"context"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
	"github.com/andreluialves/shop-orders/payments-service/internal/repository"
)

type PaymentService struct {
	paymentRepository repository.PaymentRepository
	idGenerator       repository.PaymentIDGenerator
}

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	idGenerator repository.PaymentIDGenerator,
) *PaymentService {
	return &PaymentService{
		paymentRepository: paymentRepository,
		idGenerator:       idGenerator,
	}
}

// ProcessPayment cria o registro de pagamento e decide aprová-lo ou
// recusá-lo. A regra abaixo é um placeholder — no lugar de um gateway real,
// simula recusa para valores muito altos.
func (s *PaymentService) ProcessPayment(ctx context.Context, orderID string, amount float64) (*domain.Payment, error) {
	id, err := s.idGenerator.NextPaymentID(ctx)
	if err != nil {
		return nil, err
	}

	payment, err := domain.NewPayment(id, orderID, amount)
	if err != nil {
		return nil, err
	}

	if amount > 10_000 {
		if err := payment.Decline(); err != nil {
			return nil, err
		}
	} else {
		if err := payment.Approve(); err != nil {
			return nil, err
		}
	}

	if err := s.paymentRepository.Save(payment); err != nil {
		return nil, err
	}

	return payment, nil
}
