package repository_test

import (
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/shop-orders/internal/repository"
)

func TestLoggingProductRepository_FindByID(t *testing.T) {
	t.Run("deve delegar para inner e retornar o produto", func(t *testing.T) {
		expected := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		inner := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return expected, nil
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingProductRepository(inner, log)

		result, err := repo.FindByID("P001")

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if result != expected {
			t.Error("esperava que o decorator retornasse o mesmo objeto do inner")
		}

		if errorLogged {
			t.Error("não deveria ter logado erro em caso de sucesso")
		}
	})

	t.Run("deve logar erro quando inner falha", func(t *testing.T) {
		inner := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingProductRepository(inner, log)

		_, err := repo.FindByID("P999")

		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Errorf("esperava ErrProductNotFound, recebeu %v", err)
		}

		if !errorLogged {
			t.Error("esperava que o erro fosse logado")
		}
	})
}

func TestLoggingProductRepository_Save(t *testing.T) {
	t.Run("deve delegar para inner e logar sucesso", func(t *testing.T) {
		var savedProduct *domain.Product
		var infoLogged bool

		inner := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				savedProduct = p
				return nil
			},
		}

		log := &mockLogger{
			InfoFunc: func(msg string, args ...any) {
				infoLogged = true
			},
		}

		repo := repository.NewLoggingProductRepository(inner, log)

		product := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		err := repo.Save(product)

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if savedProduct != product {
			t.Error("esperava que Save do inner fosse chamado com o mesmo produto")
		}

		if !infoLogged {
			t.Error("esperava que o sucesso fosse logado")
		}
	})

	t.Run("deve logar erro quando Save do inner falha", func(t *testing.T) {
		inner := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return errors.New("erro de conexão com banco")
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingProductRepository(inner, log)

		product := domain.RestoreProduct("P001", "Notebook", 3500, 10)

		err := repo.Save(product)

		if err == nil {
			t.Error("esperava erro propagado do Save")
		}

		if !errorLogged {
			t.Error("esperava que o erro fosse logado")
		}
	})
}

func TestLoggingProductRepository_List(t *testing.T) {
	t.Run("deve delegar para inner e retornar a lista", func(t *testing.T) {
		expected := []*domain.Product{
			domain.RestoreProduct("P001", "Notebook", 3500, 10),
			domain.RestoreProduct("P002", "Mouse", 150, 20),
		}

		inner := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return expected, nil
			},
		}

		log := &mockLogger{}

		repo := repository.NewLoggingProductRepository(inner, log)

		result, err := repo.List()

		if err != nil {
			t.Fatalf("não esperava erro, recebeu %v", err)
		}

		if len(result) != 2 {
			t.Errorf("esperava 2 produtos, recebeu %d", len(result))
		}
	})

	t.Run("deve logar erro quando List do inner falha", func(t *testing.T) {
		inner := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return nil, errors.New("erro de conexão com banco")
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingProductRepository(inner, log)

		_, err := repo.List()

		if err == nil {
			t.Error("esperava erro propagado do List")
		}

		if !errorLogged {
			t.Error("esperava que o erro fosse logado")
		}
	})
}
