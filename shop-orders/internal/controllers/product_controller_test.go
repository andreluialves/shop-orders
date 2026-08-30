package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreluialves/shop-orders/shop-orders/internal/controllers"
	"github.com/andreluialves/shop-orders/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/shop-orders/internal/service"
	"github.com/go-chi/chi/v5"
)

func TestProductController_CreateProduct(t *testing.T) {
	t.Run("deve criar produto com sucesso e retornar 201", func(t *testing.T) {
		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return nil
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		body := `{"name":"Notebook","price":3500,"quantity":10}`
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateProduct(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("esperava status 201, recebeu %d", resp.StatusCode)
		}

		var got domain.Product
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.Name != "Notebook" {
			t.Errorf("esperava name Notebook, recebeu %v", got.Name)
		}

		if got.ID != "PROD-001" {
			t.Errorf("esperava ID PROD-001, recebeu %v", got.ID)
		}
	})

	t.Run("deve retornar 400 quando body é JSON inválido", func(t *testing.T) {
		productService := service.NewProductService(&mockProductRepository{})
		controller := controllers.NewProductController(productService)

		body := `{"name": "Notebook", "price":` // JSON malformado
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateProduct(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 400 quando nome do produto é vazio", func(t *testing.T) {
		productService := service.NewProductService(&mockProductRepository{})
		controller := controllers.NewProductController(productService)

		body := `{"name":"","price":3500,"quantity":10}`
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateProduct(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 400 quando preço é inválido", func(t *testing.T) {
		productService := service.NewProductService(&mockProductRepository{})
		controller := controllers.NewProductController(productService)

		body := `{"name":"Notebook","price":-100,"quantity":10}`
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateProduct(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
		}
	})

	t.Run("deve retornar 500 quando Save falha", func(t *testing.T) {
		productRepo := &mockProductRepository{
			SaveFunc: func(p *domain.Product) error {
				return errors.New("erro de conexão com banco")
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		body := `{"name":"Notebook","price":3500,"quantity":10}`
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.CreateProduct(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", resp.StatusCode)
		}
	})
}

func TestProductController_FindAllProducts(t *testing.T) {
	t.Run("deve retornar lista de produtos com 200", func(t *testing.T) {
		expected := []*domain.Product{
			domain.RestoreProduct("PROD-001", "Notebook", 3500, 10),
			domain.RestoreProduct("PROD-002", "Mouse", 150, 20),
		}

		productRepo := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return expected, nil
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		controller.FindAllProducts(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var got []*domain.Product
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if len(got) != 2 {
			t.Errorf("esperava 2 produtos, recebeu %d", len(got))
		}
	})

	t.Run("deve retornar 500 quando List falha", func(t *testing.T) {
		productRepo := &mockProductRepository{
			ListFunc: func() ([]*domain.Product, error) {
				return nil, errors.New("erro de conexão com banco")
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		controller.FindAllProducts(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("esperava status 500, recebeu %d", resp.StatusCode)
		}
	})
}

func TestProductController_FindProductByID(t *testing.T) {
	t.Run("deve retornar produto com 200 quando encontrado", func(t *testing.T) {
		expected := domain.RestoreProduct("PROD-001", "Notebook", 3500, 10)

		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return expected, nil
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		req := httptest.NewRequest(http.MethodGet, "/products/PROD-001", nil)
		req = withURLParam(req, "id", "PROD-001")
		w := httptest.NewRecorder()

		controller.FindProductByID(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
		}

		var got domain.Product
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		if got.ID != "PROD-001" {
			t.Errorf("esperava ID PROD-001, recebeu %v", got.ID)
		}
	})

	t.Run("deve retornar 404 quando produto não existe", func(t *testing.T) {
		productRepo := &mockProductRepository{
			FindByIDFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}

		productService := service.NewProductService(productRepo)
		controller := controllers.NewProductController(productService)

		req := httptest.NewRequest(http.MethodGet, "/products/PROD-999", nil)
		req = withURLParam(req, "id", "PROD-999")
		w := httptest.NewRecorder()

		controller.FindProductByID(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("esperava status 404, recebeu %d", resp.StatusCode)
		}
	})
}

// withURLParam injeta um parâmetro de rota simulando o comportamento do chi router,
// necessário porque nos testes não passamos por um router de verdade.
func withURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
