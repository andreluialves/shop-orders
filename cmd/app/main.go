package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/repository"

	"github.com/andreluialves/shop-orders/config"
	"github.com/andreluialves/shop-orders/internal/controllers"
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

	productService := service.NewProductService(productRepository)

	// // Teste de criação de um pedido no banco de dados
	// order, err := orderService.CreateOrder(
	// 	"João Silva",
	// 	[]service.CreateOrderItem{
	// 		{
	// 			ProductID: "P001",
	// 			Quantity:  2,
	// 		},
	// 	},
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Pedido criado: %s", order.ID)

	// Controllers
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	// Cria o roteador
	router := routes.NewRouter(productController, orderController)

	log.Println("Server running on :8080")

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	// Evita warning caso os services ainda não estejam sendo utilizados
	_ = productService
	_ = orderService
}
