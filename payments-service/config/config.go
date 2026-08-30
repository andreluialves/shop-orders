package config

import (
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Database    DatabaseConfig
	RabbitMQURL string
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret               string
	AccessTokenTTL          time.Duration
	RefreshTokenIdleTTL     time.Duration
	RefreshTokenAbsoluteTTL time.Duration
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("arquivo .env nao encontrado, usando variaveis do ambiente")
	}

	database := DatabaseConfig{
		Host:     getEnv("POSTGRES_PAYMENTS_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PAYMENTS_PORT", "5432"),
		User:     getEnv("POSTGRES_PAYMENTS_USER", "shop_orders"),
		Password: getEnv("POSTGRES_PAYMENTS_PASSWORD", "shop_orders"),
		Name:     getEnv("POSTGRES_PAYMENTS_DB", "shop_orders"),
		SSLMode:  getEnv("POSTGRES_PAYMENTS_SSLMODE", "disable"),
	}

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		databaseURL = database.URL()
	}

	return Config{
		DatabaseURL: databaseURL,
		Database:    database,
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"), // novo
		Port:        getEnv("PORT", "8080"),
	}
}

func (database DatabaseConfig) URL() string {
	values := url.Values{}
	values.Set("sslmode", database.SSLMode)

	connectionURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(database.User, database.Password),
		Host:     net.JoinHostPort(database.Host, database.Port),
		Path:     "/" + database.Name,
		RawQuery: values.Encode(),
	}

	return connectionURL.String()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
