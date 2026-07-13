package dto

type CreateOrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type CreateOrderRequest struct {
	Customer string            `json:"product"`
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
