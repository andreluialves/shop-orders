package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/shop-orders/internal/logger"
	"github.com/andreluialves/shop-orders/shop-orders/internal/messaging"
	"github.com/andreluialves/shop-orders/shop-orders/internal/repository"

	"github.com/andreluialves/shop-orders/shop-orders/config"
	"github.com/andreluialves/shop-orders/shop-orders/internal/controllers"
	"github.com/andreluialves/shop-orders/shop-orders/internal/database"
	"github.com/andreluialves/shop-orders/shop-orders/internal/routes"
	"github.com/andreluialves/shop-orders/shop-orders/internal/service"
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

	appLogger := logger.NewSlogLogger()

	// Repositories
	productRepository := repository.NewPostgresProductRepository(db)
	orderRepository := repository.NewPostgresOrderRepository(db)
	customerRepository := repository.NewPostgresCustomerRepository(db)

	// Decorators com logging
	loggedOrderRepo := repository.NewLoggingOrderRepository(orderRepository, appLogger)
	loggedProductRepo := repository.NewLoggingProductRepository(productRepository, appLogger)
	loggedCustomerRepo := repository.NewLoggingCustomerRepository(customerRepository, appLogger)

	unitOfWork := repository.NewPostgresUnitOfWork(db)
	orderIDGenerator := repository.NewPostgresOrderIDGenerator(db)
	customerIDGenerator := repository.NewPostgresCustomerIDGenerator(db)

	//RabbitMQ
	rabbit, err := messaging.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer rabbit.Close()

	// Services
	orderService := service.NewOrderService(loggedProductRepo, loggedOrderRepo, unitOfWork, orderIDGenerator, rabbit)
	productService := service.NewProductService(loggedProductRepo)
	customerService := service.NewCustomerService(loggedCustomerRepo, customerIDGenerator)

	paymentHandler := messaging.NewPaymentProcessedHandler(rabbit, orderService, appLogger)
	if err := paymentHandler.Start(ctx); err != nil {
		log.Fatalf("failed to start payment handler: %v", err)
	}

	// Controllers
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)
	customerController := controllers.NewCustomerController(customerService)

	// Cria o roteador
	router := routes.NewRouter(productController, orderController, customerController)

	log.Println("Server running on :8080")

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	// Evita warning caso os services ainda não estejam sendo utilizados
	_ = productService
	_ = orderService
}
