package dto

import "github.com/andreluialves/shop-orders/shop-orders/internal/domain"

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
	CustomerID string            `json:"customer_id"`
	Items      []CreateOrderItem `json:"items"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Total      float64             `json:"total"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
}

func NewOrderResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))

	for i, item := range order.Items {
		items[i] = NewOrderItemResponse(item)
	}

	return OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Total:      order.TotalSum(),
		Status:     string(order.Status()),
		Items:      items,
	}
}
