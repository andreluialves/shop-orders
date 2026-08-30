package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
	"github.com/andreluialves/shop-orders/payments-service/internal/service"
)

func TestPaymentService_ProcessPayment(t *testing.T) {
	t.Run("deve aprovar pagamento dentro do limite", func(t *testing.T) {
		var savedPayment *domain.Payment

		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error {
				savedPayment = p
				return nil
			},
		}

		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "PAY-001", nil
			},
		}

		s := service.NewPaymentService(repo, idGen)

		payment, err := s.ProcessPayment(context.Background(), "PED-001", 500.0)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if payment.Status() != domain.PaymentStatusApproved {
			t.Errorf("esperava status APPROVED, recebeu %v", payment.Status())
		}

		if savedPayment == nil {
			t.Fatal("esperava que Save fosse chamado")
		}
	})

	t.Run("deve recusar pagamento acima do limite", func(t *testing.T) {
		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error { return nil },
		}

		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "PAY-001", nil
			},
		}

		s := service.NewPaymentService(repo, idGen)

		payment, err := s.ProcessPayment(context.Background(), "PED-001", 15000.0)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if payment.Status() != domain.PaymentStatusDeclined {
			t.Errorf("esperava status DECLINED, recebeu %v", payment.Status())
		}
	})

	t.Run("não deve criar pagamento com orderID vazio", func(t *testing.T) {
		s := service.NewPaymentService(&mockPaymentRepository{}, &mockPaymentIDGenerator{})

		_, err := s.ProcessPayment(context.Background(), "", 500.0)

		if !errors.Is(err, domain.ErrInvalidOrderID) {
			t.Errorf("esperava ErrInvalidOrderID, recebeu %v", err)
		}
	})

	t.Run("deve propagar erro quando geração de ID falha", func(t *testing.T) {
		var saveCalled bool

		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error {
				saveCalled = true
				return nil
			},
		}

		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "", errors.New("falha ao gerar ID")
			},
		}

		s := service.NewPaymentService(repo, idGen)

		_, err := s.ProcessPayment(context.Background(), "PED-001", 500.0)

		if err == nil {
			t.Error("esperava erro propagado do gerador de ID")
		}

		if saveCalled {
			t.Error("Save não deveria ser chamado quando a geração de ID falha")
		}
	})

	t.Run("deve propagar erro quando Save falha", func(t *testing.T) {
		repo := &mockPaymentRepository{
			SaveFunc: func(p *domain.Payment) error {
				return errors.New("erro de conexão com banco")
			},
		}

		idGen := &mockPaymentIDGenerator{
			NextPaymentIDFunc: func(ctx context.Context) (string, error) {
				return "PAY-001", nil
			},
		}

		s := service.NewPaymentService(repo, idGen)

		_, err := s.ProcessPayment(context.Background(), "PED-001", 500.0)

		if err == nil {
			t.Error("esperava erro propagado do Save")
		}
	})
}
