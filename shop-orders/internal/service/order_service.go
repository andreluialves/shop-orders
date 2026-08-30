package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/shop-orders/internal/events"
	"github.com/andreluialves/shop-orders/shop-orders/internal/repository"
)

type OrderService struct {
	productRepository repository.ProductRepository
	orderRepository   repository.OrderRepository
	unitOfWork        repository.UnitOfWork
	idGenerator       repository.OrderIDGenerator
	eventPublisher    events.EventPublisher
}

func NewOrderService(
	productRepository repository.ProductRepository,
	orderRepository repository.OrderRepository,
	unitOfWork repository.UnitOfWork,
	idGenerator repository.OrderIDGenerator,
	eventPublisher events.EventPublisher,
) *OrderService {
	return &OrderService{
		productRepository: productRepository,
		orderRepository:   orderRepository,
		unitOfWork:        unitOfWork,
		idGenerator:       idGenerator,
		eventPublisher:    eventPublisher,
	}
}

type CreateOrderItem struct {
	ID       string
	Quantity int
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []CreateOrderItem) (*domain.Order, error) {
	var order *domain.Order

	err := s.unitOfWork.Execute(ctx, func(repos repository.Repositories) error {
		if _, err := repos.Customer.FindByID(customerID); err != nil {
			return err
		}

		var orderItems []*domain.OrderItem

		for _, item := range items {
			product, err := repos.Product.FindByID(item.ID)
			if err != nil {
				return err
			}

			if err := product.ValidateQuantity(item.Quantity); err != nil {
				return err
			}

			orderItems = append(orderItems, domain.NewOrderItem(product, item.Quantity, product.Price))
		}

		orderID, err := s.idGenerator.NextOrderID(ctx)
		if err != nil {
			return err
		}

		newOrder, err := domain.NewOrder(orderID, customerID, orderItems)
		if err != nil {
			return err
		}

		for _, item := range newOrder.Items {
			if err := item.Product.ReduceQuantity(item.Quantity); err != nil {
				return err
			}

			if err := repos.Product.Save(item.Product); err != nil {
				return err
			}
		}

		if err := repos.Order.Save(newOrder); err != nil {
			return err
		}

		order = newOrder
		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// PayOrder inicia a saga: valida que o pedido pode ser pago e publica a
// solicitação. A confirmação real acontece de forma assíncrona, quando o
// evento payment.processed é recebido.
func (s *OrderService) PayOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if order.Status() != domain.OrderStatusPending {
		return nil, domain.ErrChangeStatusInvalid
	}

	event := events.PaymentRequested{
		OrderID: order.ID,
		Amount:  order.TotalSum(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	if err := s.eventPublisher.Publish("payment.requested", body); err != nil {
		return nil, err
	}

	return order, nil
}

// HandlePaymentResult reage ao evento payment.processed — é aqui que a
// saga se completa: confirma o pagamento (caminho feliz) ou compensa
// (desfaz a reserva de estoque e cancela o pedido) em caso de recusa.
func (s *OrderService) HandlePaymentResult(ctx context.Context, orderID string, status string) error {
	switch status {
	case "APPROVED":
		return s.confirmPayment(orderID)

	case "DECLINED":
		// Compensação da saga: aciona o CancelOrder que restaura o estoque e
		// cancela o pedido dentro de uma transação.
		_, err := s.CancelOrder(ctx, orderID)
		return err

	default:
		return fmt.Errorf("status de pagamento desconhecido: %s", status)
	}
}

func (s *OrderService) confirmPayment(orderID string) error {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if err := order.Pay(); err != nil {
		return err
	}

	return s.orderRepository.Update(order)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) (*domain.Order, error) {
	var order *domain.Order

	err := s.unitOfWork.Execute(ctx, func(repos repository.Repositories) error {
		o, err := repos.Order.FindByID(id)
		if err != nil {
			return err
		}

		if err := o.Cancel(); err != nil {
			return err
		}

		for _, item := range o.Items {
			product, err := repos.Product.FindByID(item.Product.ID)
			if err != nil {
				return err
			}

			if err := product.RestoreQuantity(item.Quantity); err != nil {
				return err
			}

			if err := repos.Product.Save(product); err != nil {
				return err
			}
		}

		if err := repos.Order.Update(o); err != nil {
			return err
		}

		order = o
		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) FindOrderByID(id string) (*domain.Order, error) {
	return s.orderRepository.FindByID(id)
}

func (s *OrderService) ListOrders(limit, offset int) ([]*domain.Order, int, error) {
	return s.orderRepository.List(limit, offset)
}

type OrderFilter func(*domain.Order) bool

func (s *OrderService) FilterOrders(filter OrderFilter, limit, offset int) ([]*domain.Order, int, error) {

	orders, total, err := s.orderRepository.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	filteredOrders := make([]*domain.Order, 0)

	for _, order := range orders {
		if filter(order) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders, total, nil
}
