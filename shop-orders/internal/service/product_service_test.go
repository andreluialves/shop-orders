package service_test

import (
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/service"
)

func TestProductService_CreateProduct(t *testing.T) {
	t.Run("deve criar produto e atribuir ID gerado", func(t *testing.T) {
		var savedProduct *domain.Product

		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				savedProduct = p
				return nil
			},
		}

		s := service.NewProductService(productRepo)

		product := &domain.Product{
			Name:     "Notebook",
			Price:    3500,
			Quantity: 10,
		}

		err := s.CreateProduct(product)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if product.ID != "PROD-001" {
			t.Errorf("esperava ID PROD-001, recebeu %v", product.ID)
		}

		if savedProduct == nil {
			t.Fatal("esperava que Save fosse chamado, mas não foi")
		}

		if savedProduct.ID != "PROD-001" {
			t.Errorf("esperava que Save recebesse produto com ID PROD-001, recebeu %v", savedProduct.ID)
		}
	})

	t.Run("deve incrementar o ID a cada produto criado", func(t *testing.T) {
		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		s := service.NewProductService(productRepo)

		product1 := &domain.Product{Name: "Notebook", Price: 3500, Quantity: 10}
		product2 := &domain.Product{Name: "Mouse", Price: 150, Quantity: 20}

		_ = s.CreateProduct(product1)
		_ = s.CreateProduct(product2)

		if product1.ID != "PROD-001" {
			t.Errorf("esperava ID PROD-001 no primeiro produto, recebeu %v", product1.ID)
		}

		if product2.ID != "PROD-002" {
			t.Errorf("esperava ID PROD-002 no segundo produto, recebeu %v", product2.ID)
		}
	})

	t.Run("deve propagar erro quando Save falha", func(t *testing.T) {
		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return errors.New("erro de conexão com banco")
			},
		}

		s := service.NewProductService(productRepo)

		product := &domain.Product{Name: "Notebook", Price: 3500, Quantity: 10}

		err := s.CreateProduct(product)

		if err == nil {
			t.Error("esperava erro propagado do Save, recebeu nil")
		}
	})

	t.Run("CreateProduct do service não valida — depende do domain.NewProduct ter sido chamado antes", func(t *testing.T) {
		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		s := service.NewProductService(productRepo)

		// Simula um produto inválido chegando direto no service,
		// sem passar por domain.NewProduct (ex: chamada interna, sem HTTP)
		invalidProduct := &domain.Product{Name: "", Price: -100, Quantity: -5}

		err := s.CreateProduct(invalidProduct)

		if err != nil {
			t.Errorf("comportamento atual esperado: Save não é bloqueado por dados inválidos; recebeu erro %v", err)
		}
	})
}

func TestProductService_FindByID(t *testing.T) {
	t.Run("deve retornar produto quando encontrado", func(t *testing.T) {
		expected := domain.RestoreProduct("PROD-001", "Notebook", 3500, 10)

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return expected, nil
			},
		}

		s := service.NewProductService(productRepo)

		result, err := s.FindByID("PROD-001")

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if result != expected {
			t.Errorf("esperava o produto retornado pelo repository, recebeu outro")
		}
	})

	t.Run("deve retornar erro quando produto não é encontrado", func(t *testing.T) {
		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}

		s := service.NewProductService(productRepo)

		_, err := s.FindByID("PROD-999")

		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Errorf("esperava ErrProductNotFound, recebeu %v", err)
		}
	})
}

func TestProductService_List(t *testing.T) {
	t.Run("deve retornar lista de produtos", func(t *testing.T) {
		expected := []*domain.Product{
			domain.RestoreProduct("PROD-001", "Notebook", 3500, 10),
			domain.RestoreProduct("PROD-002", "Mouse", 150, 20),
		}

		productRepo := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return expected, nil
			},
		}

		s := service.NewProductService(productRepo)

		result, err := s.List()

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if len(result) != 2 {
			t.Errorf("esperava 2 produtos, recebeu %v", len(result))
		}
	})

	t.Run("deve propagar erro quando List falha", func(t *testing.T) {
		productRepo := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return nil, errors.New("erro de conexão com banco")
			},
		}

		s := service.NewProductService(productRepo)

		_, err := s.List()

		if err == nil {
			t.Error("esperava erro propagado do List, recebeu nil")
		}
	})
}
