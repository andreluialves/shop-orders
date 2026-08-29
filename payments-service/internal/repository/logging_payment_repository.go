package repository

import (
	"github.com/andreluialves/shop-orders/payments-service/internal/domain"
	"github.com/andreluialves/shop-orders/payments-service/internal/logger"
)

type LoggingPaymentRepository struct {
	inner  PaymentRepository
	logger logger.Logger
}

func NewLoggingPaymentRepository(inner PaymentRepository, logger logger.Logger) *LoggingPaymentRepository {
	return &LoggingPaymentRepository{inner: inner, logger: logger}
}

func (l *LoggingPaymentRepository) FindByID(id string) (*domain.Payment, error) {
	l.logger.Debug("buscando pagamento", "payment_id", id)

	payment, err := l.inner.FindByID(id)
	if err != nil {
		l.logger.Error("falha ao buscar pagamento", "payment_id", id, "error", err)
		return nil, err
	}

	return payment, nil
}

func (l *LoggingPaymentRepository) Save(payment *domain.Payment) error {
	l.logger.Info("salvando pagamento", "payment_id", payment.ID)

	if err := l.inner.Save(payment); err != nil {
		l.logger.Error("falha ao salvar pagamento", "payment_id", payment.ID, "error", err)
		return err
	}

	return nil
}

func (l *LoggingPaymentRepository) List() ([]*domain.Payment, error) {
	l.logger.Debug("listando pagamentos")

	payments, err := l.inner.List()
	if err != nil {
		l.logger.Error("falha ao listar pagamentos", "error", err)
		return nil, err
	}

	return payments, nil
}
