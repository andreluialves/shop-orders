package service

import (
	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
)

type ProductService struct {
	productRepository repository.ProductRepository
}

func NewProductService(
	productRepository repository.ProductRepository,
) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	return s.productRepository.Save(product)
}

func (s *ProductService) FindByID(id string) (*domain.Product, error) {
	return s.productRepository.FindByID(id)
}

func (s *ProductService) List() ([]*domain.Product, error) {
	return s.productRepository.List()
}
