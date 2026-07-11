package service

import (
	"fmt"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
)

type OrderService struct {
	productRepository repository.ProductRepository
	orderRepository   repository.OrderRepository
	nextOrderID       int
}

func NewOrderService(
	productRepository repository.ProductRepository,
	orderRepository repository.OrderRepository,
) *OrderService {
	return &OrderService{
		productRepository: productRepository,
		orderRepository:   orderRepository,
		nextOrderID:       0,
	}
}

func (s *OrderService) generateOrderID() string {
	s.nextOrderID++
	return fmt.Sprintf("PED-%03d", s.nextOrderID)
}

type CreateOrderItem struct {
	ProductID string
	Quantity  int
}

func (s *OrderService) CreateOrder(customer string, items []CreateOrderItem) (*domain.Order, error) {
	var orderItems []*domain.OrderItem

	for _, item := range items {
		product, err := s.productRepository.FindByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		orderItem := domain.NewOrderItem(product, item.Quantity, product.Price)
		orderItems = append(orderItems, orderItem)
	}

	order, err := domain.NewOrder(s.generateOrderID(), customer, orderItems)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepository.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) PayOrder(id string) error {
	order, err := s.FindOrderByID(id)
	if err != nil {
		return err
	}
	order.Pay()
	return s.orderRepository.Save(order)
}

func (s *OrderService) CancelOrder(id string) error {

	order, err := s.FindOrderByID(id)

	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	for _, item := range order.Items {

		product, err := s.productRepository.FindByID(item.Product.ID)

		if err != nil {
			return err
		}

		if err := product.RestoreQuantity(item.Quantity); err != nil {
			return err
		}

		if err := s.productRepository.Save(product); err != nil {
			return err
		}
	}

	return s.orderRepository.Save(order)
}

func (s *OrderService) FindOrderByID(id string) (*domain.Order, error) {
	return s.orderRepository.FindByID(id)
}

func (s *OrderService) ListOrders() ([]*domain.Order, error) {
	return s.orderRepository.List()
}
