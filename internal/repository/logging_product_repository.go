package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/logger"
)

type LoggingProductRepository struct {
	inner  ProductRepository
	logger logger.Logger
}

func NewLoggingProductRepository(inner ProductRepository, logger logger.Logger) *LoggingProductRepository {
	return &LoggingProductRepository{inner: inner, logger: logger}
}

func (l *LoggingProductRepository) FindByID(id string) (*domain.Product, error) {
	l.logger.Debug("buscando produto", "product_id", id)

	product, err := l.inner.FindByID(id)
	if err != nil {
		l.logger.Error("falha ao buscar produto", "product_id", id, "error", err)
		return nil, err
	}

	return product, nil
}

func (l *LoggingProductRepository) Save(product *domain.Product) error {
	l.logger.Info("salvando produto", "product_id", product.ID, "quantity", product.Quantity)

	if err := l.inner.Save(product); err != nil {
		l.logger.Error("falha ao salvar produto", "product_id", product.ID, "error", err)
		return err
	}

	return nil
}

func (l *LoggingProductRepository) List() ([]*domain.Product, error) {
	l.logger.Debug("listando produtos")

	products, err := l.inner.List()
	if err != nil {
		l.logger.Error("falha ao listar produtos", "error", err)
		return nil, err
	}

	return products, nil
}
