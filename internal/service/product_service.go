package service

import (
	"fmt"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
)

type ProductService struct {
	productRepository repository.ProductRepository
	nextProductID     int
}

func (s *ProductService) generateProductID() string {
	s.nextProductID++
	return fmt.Sprintf("PROD-%03d", s.nextProductID)
}

func NewProductService(
	productRepository repository.ProductRepository,
) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {

	product.ID = s.generateProductID()

	return s.productRepository.Save(product)
}

func (s *ProductService) FindByID(id string) (*domain.Product, error) {
	return s.productRepository.FindByID(id)
}

func (s *ProductService) List() ([]*domain.Product, error) {
	return s.productRepository.List()
}
