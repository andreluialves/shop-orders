package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/service"

	"github.com/andreluialves/shop-orders/internal/dto"
	"github.com/go-chi/chi/v5"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := domain.NewProduct(
		request.Name,
		request.Price,
		request.Quantity,
	)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := c.productService.CreateProduct(product); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (c *ProductController) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productService.List()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (c *ProductController) FindProductByID(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	product, err := c.productService.FindByID(vars)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
