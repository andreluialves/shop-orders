package repository

import (
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/logger"
)

type LoggingOrderRepository struct {
	inner  OrderRepository // a implementação real "por dentro" (ex: PostgresOrderRepository)
	logger logger.Logger
}

func NewLoggingOrderRepository(inner OrderRepository, logger logger.Logger) *LoggingOrderRepository {
	return &LoggingOrderRepository{
		inner:  inner,
		logger: logger,
	}
}

func (l *LoggingOrderRepository) FindByID(id string) (*domain.Order, error) {
	l.logger.Debug("buscando pedido", "order_id", id)

	order, err := l.inner.FindByID(id)
	if err != nil {
		l.logger.Error("falha ao buscar pedido", "order_id", id, "error", err)
		return nil, err
	}

	l.logger.Debug("pedido encontrado", "order_id", id)
	return order, nil
}

func (l *LoggingOrderRepository) Save(order *domain.Order) error {
	l.logger.Info("salvando novo pedido", "order_id", order.ID, "customer", order.Customer)

	if err := l.inner.Save(order); err != nil {
		l.logger.Error("falha ao salvar pedido", "order_id", order.ID, "error", err)
		return err
	}

	l.logger.Info("pedido salvo com sucesso", "order_id", order.ID)
	return nil
}

func (l *LoggingOrderRepository) Update(order *domain.Order) error {
	l.logger.Info("atualizando pedido", "order_id", order.ID, "status", order.Status())

	if err := l.inner.Update(order); err != nil {
		l.logger.Error("falha ao atualizar pedido", "order_id", order.ID, "error", err)
		return err
	}

	l.logger.Info("pedido atualizado com sucesso", "order_id", order.ID)
	return nil
}

func (l *LoggingOrderRepository) List(limit, offset int) ([]*domain.Order, int, error) {
	l.logger.Debug("listando pedidos", "limit", limit, "offset", offset)

	orders, total, err := l.inner.List(limit, offset)
	if err != nil {
		l.logger.Error("falha ao listar pedidos", "error", err)
		return nil, 0, err
	}

	l.logger.Debug("pedidos listados", "total", total, "retornados", len(orders))
	return orders, total, nil
}
