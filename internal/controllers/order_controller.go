package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/internal/dto"
	"github.com/andreluialves/shop-orders/internal/service"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController(
	orderService *service.OrderService,
) *OrderController {

	return &OrderController{
		orderService: orderService,
	}
}

func (oc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var request dto.CreateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]service.CreateOrderItem, 0)

	for _, item := range request.Items {
		items = append(items, service.CreateOrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := oc.orderService.CreateOrder(
		request.Customer,
		items,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) FindAllOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := oc.orderService.ListOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) FindOrderByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	order, err := oc.orderService.FindOrderByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(order)
}
