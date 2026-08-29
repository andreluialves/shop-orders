package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/logger"
)

type LoggingCustomerRepository struct {
	inner  CustomerRepository
	logger logger.Logger
}

func NewLoggingCustomerRepository(inner CustomerRepository, logger logger.Logger) *LoggingCustomerRepository {
	return &LoggingCustomerRepository{inner: inner, logger: logger}
}

func (l *LoggingCustomerRepository) FindByID(id string) (*domain.Customer, error) {
	l.logger.Debug("buscando cliente", "customer_id", id)

	customer, err := l.inner.FindByID(id)
	if err != nil {
		l.logger.Error("falha ao buscar cliente", "customer_id", id, "error", err)
		return nil, err
	}

	return customer, nil
}

func (l *LoggingCustomerRepository) Save(customer *domain.Customer) error {
	l.logger.Info("salvando cliente", "customer_id", customer.ID)

	if err := l.inner.Save(customer); err != nil {
		l.logger.Error("falha ao salvar cliente", "customer_id", customer.ID, "error", err)
		return err
	}

	return nil
}

func (l *LoggingCustomerRepository) List() ([]*domain.Customer, error) {
	l.logger.Debug("listando clientes")

	customers, err := l.inner.List()
	if err != nil {
		l.logger.Error("falha ao listar clientes", "error", err)
		return nil, err
	}

	return customers, nil
}
