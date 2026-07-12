package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/config"
	"github.com/andreluialves/shop-orders/internal/database"
	"github.com/andreluialves/shop-orders/internal/routes"
)

func main() {

	// Carrega as configurações da aplicação
	cfg := config.Load()

	// Contexto utilizado para inicialização
	ctx := context.Background()

	// Cria o pool de conexões com o PostgreSQL
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	// Cria o roteador
	router := routes.NewRouter()

	log.Println("Server running on :8080")

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
