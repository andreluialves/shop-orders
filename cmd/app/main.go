package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/routes"
)

func main() {
	// carregar configuração

	// abrir conexão

	// criar repositories

	// criar services

	// criar controllers

	// registrar rotas

	// iniciar servidor
	router := routes.NewRouter()

	fmt.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
