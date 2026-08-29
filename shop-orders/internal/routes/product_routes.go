package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreluialves/shop-orders/internal/controllers"
)

func ProductRoutes(r chi.Router, productController *controllers.ProductController) {
	r.Post("/products", productController.CreateProduct)
	r.Get("/products", productController.FindAllProducts)
	r.Get("/products/{id}", productController.FindProductByID)
}
