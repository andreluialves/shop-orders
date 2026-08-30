package domain_test

import (
	"testing"

	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
)

func TestNewPayment(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		amount  float64
		wantErr error
	}{
		{"deve criar pagamento válido", "PED-001", 100.0, nil},
		{"não deve criar sem orderID", "", 100.0, domain.ErrInvalidOrderID},
		{"não deve criar com amount zero", "PED-001", 0, domain.ErrInvalidAmount},
		{"não deve criar com amount negativo", "PED-001", -10, domain.ErrInvalidAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := domain.NewPayment("PAY-001", tt.orderID, tt.amount)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				if payment != nil {
					t.Error("esperava payment nil em caso de erro")
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if payment.Status() != domain.PaymentStatusPending {
				t.Errorf("esperava status PENDING, recebeu %v", payment.Status())
			}
		})
	}
}

func TestPayment_Approve(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus domain.PaymentStatus
		wantErr       error
	}{
		{"deve aprovar pagamento pendente", domain.PaymentStatusPending, nil},
		{"não deve aprovar pagamento já aprovado", domain.PaymentStatusApproved, domain.ErrPaymentAlreadyDone},
		{"não deve aprovar pagamento já recusado", domain.PaymentStatusDeclined, domain.ErrPaymentAlreadyDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := domain.RestorePayment("PAY-001", "PED-001", 100.0, tt.initialStatus)

			err := payment.Approve()

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if payment.Status() != domain.PaymentStatusApproved {
				t.Errorf("esperava status APPROVED, recebeu %v", payment.Status())
			}
		})
	}
}

func TestPayment_Decline(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus domain.PaymentStatus
		wantErr       error
	}{
		{"deve recusar pagamento pendente", domain.PaymentStatusPending, nil},
		{"não deve recusar pagamento já aprovado", domain.PaymentStatusApproved, domain.ErrPaymentAlreadyDone},
		{"não deve recusar pagamento já recusado", domain.PaymentStatusDeclined, domain.ErrPaymentAlreadyDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := domain.RestorePayment("PAY-001", "PED-001", 100.0, tt.initialStatus)

			err := payment.Decline()

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("não esperava erro, recebeu %v", err)
			}

			if payment.Status() != domain.PaymentStatusDeclined {
				t.Errorf("esperava status DECLINED, recebeu %v", payment.Status())
			}
		})
	}
}
