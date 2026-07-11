package repository

import (
	"fmt"

	"github.com/andreluialves/shop-orders/internal/domain"
)

type MemoryProductRepository struct {
	products map[string]*domain.Product
}

func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (r *MemoryProductRepository) Save(product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) FindByID(id string) (*domain.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("product %s: %w", id, domain.ErrProductNotFound)
	}

	return product, nil
}

func (r *MemoryProductRepository) List() ([]*domain.Product, error) {
	products := make([]*domain.Product, 0, len(r.products))

	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}
