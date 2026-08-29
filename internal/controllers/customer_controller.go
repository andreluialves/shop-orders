package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/dto"
	"github.com/andreluialves/shop-orders/internal/service"
	"github.com/go-chi/chi/v5"
)

type CustomerController struct {
	customerService *service.CustomerService
}

func NewCustomerController(customerService *service.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: customerService,
	}
}

func (cc *CustomerController) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	var request dto.CreateCustomerRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customer, err := domain.NewCustomer(
		request.Name,
		request.Email,
		request.Address,
		request.Phone,
	)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := cc.customerService.CreateCustomer(r.Context(), customer); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(customer); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (c *CustomerController) FindAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := c.customerService.List()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(customers); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (c *CustomerController) FindCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	customer, err := c.customerService.FindByID(vars)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(customer); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
