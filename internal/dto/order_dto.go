package dto

import "github.com/andreluialves/shop-orders/internal/domain"

type CreateOrderItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderItemRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderItemResponse struct {
	ID       string  `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

func NewOrderItemResponse(item *domain.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:       item.Product.ID,
		Quantity: item.Quantity,
		Price:    item.Price,
		Subtotal: item.Price * float64(item.Quantity),
	}
}

type CreateOrderRequest struct {
	Customer string            `json:"customer"`
	Items    []CreateOrderItem `json:"items"`
	Status   string            `json:"status"`
}

type OrderResponse struct {
	ID       string              `json:"id"`
	Customer string              `json:"customer"`
	Total    float64             `json:"total"`
	Status   string              `json:"status"`
	Items    []OrderItemResponse `json:"items"`
}

func NewOrderResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))

	for i, item := range order.Items {
		items[i] = NewOrderItemResponse(item)
	}

	return OrderResponse{
		ID:       order.ID,
		Customer: order.Customer,
		Total:    order.TotalSum(),
		Status:   string(order.Status()),
		Items:    items,
	}
}
