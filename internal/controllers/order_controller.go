package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/internal/dto"
	"github.com/andreluialves/shop-orders/internal/pagination"
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
	log.Printf("Request: %+v\n", request)

	for i, item := range request.Items {
		log.Printf("Item %d: ID=%s Quantity=%d", i, item.ID, item.Quantity)
	}

	items := make([]service.CreateOrderItem, 0)

	for _, item := range request.Items {
		items = append(items, service.CreateOrderItem{
			ID:       item.ID,
			Quantity: item.Quantity,
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
	p := pagination.ParseFromRequest(r)

	orders, total, err := oc.orderService.ListOrders(p.Limit, p.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		data[i] = dto.NewOrderResponse(o)
	}

	response := struct {
		Data   []dto.OrderResponse `json:"data"`
		Total  int                 `json:"total"`
		Limit  int                 `json:"limit"`
		Offset int                 `json:"offset"`
	}{
		Data:   data,
		Total:  total,
		Limit:  p.Limit,
		Offset: p.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (oc *OrderController) FindOrderByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	order, err := oc.orderService.FindOrderByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.NewOrderResponse(order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (oc *OrderController) PayOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "O ID do pedido é obrigatório.", http.StatusBadRequest)
		return
	}

	order, err := oc.orderService.PayOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.NewOrderResponse(order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (oc *OrderController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "O ID do pedido é obrigatório.", http.StatusBadRequest)
		return
	}

	order, err := oc.orderService.CancelOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.NewOrderResponse(order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
