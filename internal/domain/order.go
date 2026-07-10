package domain

import "strings"

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	Product  *Product
	Quantity int
	Price    float64
}

func NewOrderItem(product *Product, quantity int, price float64) *OrderItem {
	return &OrderItem{
		Product:  product,
		Quantity: quantity,
		Price:    price,
	}
}

type Order struct {
	ID       string
	Customer string
	Items    []*OrderItem
	status   OrderStatus
}

func NewOrder(id string, customer string) (*Order, error) {
	order := Order{
		ID:       id,
		Customer: customer,
		Items:    []*OrderItem{},
		status:   OrderStatusPending,
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (o Order) Validate() error {
	if strings.TrimSpace(o.Customer) == "" {
		return ErrInvalidCustomer
	}

	if len(o.Items) == 0 {
		return ErrEmptyOrder
	}

	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (i OrderItem) Validate() error {
	if i.Product == nil {
		return ErrProductNotFound
	}

	if i.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	return nil
}

func TotalSum(orders []*Order) float64 {
	var total float64
	for _, order := range orders {
		for _, item := range order.Items {
			total += item.Price * float64(item.Quantity)
		}
	}
	return total
}

func PayOrder(order *Order) {
	order.status = OrderStatusPaid
}

func CancelOrder(order *Order) {
	order.status = OrderStatusCanceled
}
