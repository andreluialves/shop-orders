package dto

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

type CreateOrderRequest struct {
	Customer string            `json:"customer"`
	Items    []CreateOrderItem `json:"items"`
	Status   string            `json:"status"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Total      float64             `json:"total"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
}
