package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis do ambiente")
	}

	host := os.Getenv("POSTGRES_ORDERS_HOST")
	port := os.Getenv("POSTGRES_ORDERS_PORT")
	user := os.Getenv("POSTGRES_ORDERS_USER")
	password := os.Getenv("POSTGRES_ORDERS_PASSWORD")
	dbname := os.Getenv("POSTGRES_ORDERS_DB")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Variáveis de conexão do PostgreSQL não configuradas")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	content, err := os.ReadFile("seeds/100001_seed_products.up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Seed executado com sucesso")
}
