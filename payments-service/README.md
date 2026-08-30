# Payments Service

Serviço responsável por processar pagamentos de pedidos, de forma assíncrona.

## Responsabilidade

- Recebe solicitações de pagamento via evento `payment.requested`
- Decide aprovação/recusa e persiste o registro do pagamento
- Publica o resultado via evento `payment.processed`

## Estrutura do Serviço

```text

├── cmd
│   └── app
│       └── main.go
├── config
│   └── config.go
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   └── postgres.go
│   ├── domain
│   │   ├── payment.go
│   │   └── payment_test.go
│   ├── dto
│   │   └── payment_dto.go
│   ├── logger
│   │   ├── logger.go
│   │   └── slog_logger.go
│   ├── messaging
│   │   ├── events.go
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── mocks_test.go
│   │   └── rabbitmq.go
│   ├── repository
│   │   ├── dbtx.go
│   │   ├── logging_payment_repository.go
│   │   ├── payment_id_generator.go
│   │   ├── payment_repository.go
│   │   ├── postgres_payment_id_generator.go
│   │   └── postgres_payment_repository.go
│   └── service
│       ├── mocks_test.go
│       ├── payment_service.go
│       └── payment_service_test.go
├── Makefile
├── migrations
│   ├── 000001_create_payments.down.sql
│   └── 000001_create_payments.up.sql
└── README.md
```

## Configuração

Copie `.env.example` para `.env` e ajuste conforme necessário.

## Rodando

```bash
make run
```

## Migrations

```bash
make migrate-up
```