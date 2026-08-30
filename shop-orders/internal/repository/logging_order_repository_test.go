package repository_test

import (
	"errors"
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/shop-orders/internal/repository"
)

func TestLoggingOrderRepository_FindByID(t *testing.T) {
	t.Run("deve delegar para inner e retornar o pedido", func(t *testing.T) {
		expected := domain.RestoreOrder("PED-001", "João Silva", domain.OrderStatusPending)

		inner := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return expected, nil
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingOrderRepository(inner, log)

		result, err := repo.FindByID("PED-001")

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
		inner := &mockOrderRepository{
			FindByIDFunc: func(id string) (*domain.Order, error) {
				return nil, domain.ErrOrderNotFound
			},
		}

		var errorLogged bool
		log := &mockLogger{
			ErrorFunc: func(msg string, args ...any) {
				errorLogged = true
			},
		}

		repo := repository.NewLoggingOrderRepository(inner, log)

		_, err := repo.FindByID("PED-999")

		if !errors.Is(err, domain.ErrOrderNotFound) {
			t.Errorf("esperava ErrOrderNotFound, recebeu %v", err)
		}

		if !errorLogged {
			t.Error("esperava que o erro fosse logado")
		}
	})
}
