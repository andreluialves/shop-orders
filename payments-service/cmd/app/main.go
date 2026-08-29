package app

import (
	"context"
	"log"

	"github.com/andreluialves/shop-orders/payments-service/internal/logger"
	"github.com/andreluialves/shop-orders/payments-service/internal/messaging"
	"github.com/andreluialves/shop-orders/payments-service/internal/repository"

	"github.com/andreluialves/shop-orders/payments-service/config"
	"github.com/andreluialves/shop-orders/payments-service/internal/database"
	"github.com/andreluialves/shop-orders/payments-service/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	rabbit, err := messaging.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer rabbit.Close()

	appLogger := logger.NewSlogLogger()

	paymentRepo := repository.NewPostgresPaymentRepository(db)
	idGenerator := repository.NewPostgresPaymentIDGenerator(db)

	loggedPaymentRepo := repository.NewLoggingPaymentRepository(paymentRepo, appLogger)
	paymentService := service.NewPaymentService(loggedPaymentRepo, idGenerator)

	handler := messaging.NewPaymentRequestHandler(rabbit, paymentService, appLogger)
	if err := handler.Start(ctx); err != nil {
		log.Fatalf("failed to start payment handler: %v", err)
	}

	log.Println("Payments service listening for payment.requested events...")
	select {} // mantém o processo vivo
}
