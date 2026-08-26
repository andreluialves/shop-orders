package service

import (
	"context"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
)

type OrderService struct {
	productRepository repository.ProductRepository
	orderRepository   repository.OrderRepository
	unitOfWork        repository.UnitOfWork
	idGenerator       repository.OrderIDGenerator
}

func NewOrderService(
	productRepository repository.ProductRepository,
	orderRepository repository.OrderRepository,
	unitOfWork repository.UnitOfWork,
	idGenerator repository.OrderIDGenerator,
) *OrderService {
	return &OrderService{
		productRepository: productRepository,
		orderRepository:   orderRepository,
		unitOfWork:        unitOfWork,
		idGenerator:       idGenerator,
	}
}

type CreateOrderItem struct {
	ID       string
	Quantity int
}

func (s *OrderService) CreateOrder(ctx context.Context, customer string, items []CreateOrderItem) (*domain.Order, error) {
	var order *domain.Order

	err := s.unitOfWork.Execute(ctx, func(repos repository.Repositories) error {
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

		newOrder, err := domain.NewOrder(orderID, customer, orderItems)
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

func (s *OrderService) PayOrder(id string) (*domain.Order, error) {
	order, err := s.orderRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := order.Pay(); err != nil {
		return nil, err
	}

	if err := s.orderRepository.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

// func (s *OrderService) CancelOrder(id string) (*domain.Order, error) {
// 	order, err := s.FindOrderByID(id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := order.Cancel(); err != nil {
// 		return nil, err
// 	}

// 	for _, item := range order.Items {

// 		product, err := s.productRepository.FindByID(item.Product.ID)

// 		if err != nil {
// 			return nil, err
// 		}

// 		if err := product.RestoreQuantity(item.Quantity); err != nil {
// 			return nil, err
// 		}

// 		if err := s.productRepository.Save(product); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if err := s.orderRepository.Update(order); err != nil {
// 		return nil, err
// 	}

// 	return order, nil
// }

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
