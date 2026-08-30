include .env
export

DATABASE_URL = postgres://$(POSTGRES_ORDERS_USER):$(POSTGRES_ORDERS_PASSWORD)@$(POSTGRES_ORDERS_HOST):$(POSTGRES_ORDERS_PORT)/$(POSTGRES_ORDERS_DB)?sslmode=$(POSTGRES_ORDERS_SSLMODE)
MIGRATIONS_PATH = migrations
APP_MAIN = cmd/app/main.go
SEED_MAIN = cmd/seed/main.go

.PHONY: help build run run-seed test test-verbose coverage coverage-html \
        migrate-up migrate-down migrate-version migrate-create \
        docker-up docker-down docker-logs \
        fmt vet tidy clean

## help: lista todos os comandos disponíveis
help:
	@echo "Comandos disponíveis:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: compila a aplicação
build:
	go build ./...

## run: roda a aplicação (cmd/app)
run:
	go run $(APP_MAIN)

## run-seed: roda o seed de dados (cmd/seed)
run-seed:
	go run $(SEED_MAIN)

## test: roda todos os testes
test:
	go test ./...

## test-verbose: roda todos os testes com saída detalhada
test-verbose:
	go test ./... -v

## coverage: roda os testes mostrando percentual de cobertura por pacote
coverage:
	go test ./... -cover

## coverage-html: gera relatório de cobertura em HTML e abre no navegador
coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

## migrate-up: aplica todas as migrations pendentes
migrate-up:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) up

## migrate-down: desfaz a última migration aplicada
migrate-down:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) down 1

## migrate-version: mostra a versão atual da migration aplicada
migrate-version:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) version

## migrate-create: cria uma nova migration (uso: make migrate-create name=nome_da_migration)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

## docker-up: sobe os containers (Postgres, etc)
docker-up:
	docker-compose up -d

## docker-down: derruba os containers
docker-down:
	docker-compose down

## docker-logs: mostra os logs dos containers em tempo real
docker-logs:
	docker-compose logs -f

## fmt: formata todo o código Go do projeto
fmt:
	go fmt ./...

## vet: roda análise estática do Go (detecta erros comuns)
vet:
	go vet ./...

## tidy: organiza e limpa as dependências do go.mod/go.sum
tidy:
	go mod tidy

## clean: remove arquivos gerados (binários, cobertura)
clean:
	rm -f coverage.out
	go clean