package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"

	"github.com/andreluialves/shop-orders/config"
	"github.com/andreluialves/shop-orders/internal/database"
	"github.com/andreluialves/shop-orders/internal/routes"
	"github.com/andreluialves/shop-orders/internal/service"
)

func main() {

	// Carrega as configurações da aplicação
	cfg := config.Load()

	// Contexto utilizado para inicialização
	ctx := context.Background()

	// Cria o pool de conexões com o PostgreSQL
	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// Repositories
	productRepository := repository.NewPostgresProductRepository(db)
	orderRepository := repository.NewPostgresOrderRepository(db)

	// Services
	orderService := service.NewOrderService(
		productRepository,
		orderRepository,
	)

	productService := service.NewProductService(
		productRepository,
	)

	// Teste criação de um produto no banco de dados
	product, err := domain.NewProduct(
		"P001",
		"Notebook",
		3500.00,
		10,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(product); err != nil {
		log.Fatal(err)
	}

	log.Printf("Produto criado: %s", product.ID)

	// Teste de criação de um pedido no banco de dados
	order, err := orderService.CreateOrder(
		"João Silva",
		[]service.CreateOrderItem{
			{
				ProductID: "P001",
				Quantity:  2,
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Pedido criado: %s", order.ID)

	// Cria o roteador
	router := routes.NewRouter()

	// Temporariamente ainda não existem controllers.
	// Quando eles forem implementados, serão registrados aqui.

	log.Println("Server running on :8080")

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	// Evita warning caso os services ainda não estejam sendo utilizados
	_ = productService
	_ = orderService
}
