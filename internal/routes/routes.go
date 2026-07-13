package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/internal/controllers"
)

func NewRouter(
	productController *controllers.ProductController,
	orderController *controllers.OrderController,
) *chi.Mux {

	r := chi.NewRouter()

	ProductRoutes(r, productController)
	OrderRoutes(r, orderController)

	return r
}
