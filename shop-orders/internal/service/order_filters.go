package service

import "github.com/andreluialves/shop-orders/internal/domain"

func PaidOrders() OrderFilter {
	return func(order *domain.Order) bool {
		return order.Status() == domain.OrderStatusPaid
	}
}

func PendingOrders() OrderFilter {
	return func(order *domain.Order) bool {
		return order.Status() == domain.OrderStatusPending
	}
}

func OrdersAboveValue(value float64) OrderFilter {
	return func(order *domain.Order) bool {
		return order.TotalSum() > value
	}
}
