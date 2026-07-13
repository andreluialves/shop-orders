package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/internal/controllers"
)

func OrderRoutes(r chi.Router, orderController *controllers.OrderController) {
	r.Post("/orders", orderController.CreateOrder)
	r.Get("/orders", orderController.FindAllOrders)
	r.Get("/orders/{id}", orderController.FindOrderByID)
}
