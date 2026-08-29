package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/shop-orders/internal/controllers"
)

func CustomerRoutes(r chi.Router, customerController *controllers.CustomerController) {
	r.Post("/customers", customerController.CreateCustomer)
	r.Get("/customers", customerController.FindAllCustomers)
	r.Get("/customers/{id}", customerController.FindCustomerByID)
}
