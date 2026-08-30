# Shop Orders

Projeto desenvolvido em Go para estudo de arquitetura em camadas, boas práticas de desenvolvimento, criação de APIs REST e aplicação de regras de negócio no gerenciamento de produtos e pedidos.

## Funcionalidades

### Fase 1 - Domínio e regras de negócio e persistẽncia em memória

* Cadastro de produtos
* Criação de pedidos
* Validação de clientes e itens
* Controle de estoque
* Pagamento de pedidos
* Cancelamento de pedidos com restauração de estoque
* Listagem de pedidos
* Filtros de pedidos por:

  * Status (pendentes e pagos)
  * Valor mínimo do pedido

### Fase 2 - API e persistência em banco de dados (em andamento)

* Implementação da API REST de produtos
* Implementação da API REST de pedidos 
* Implementação da API REST de clientes
* Persistência de dados utilizando PostgreSQL
* Ambiente de banco de dados configurado com Docker

### Fase 3 - Mensageria e Saga

* Criação de pedidos com transação
* Cadastrar clientes
* Listar clientes
* Criar pedido para um cliente
* Buscar cliente por id
* Serviço de pagamentos "Payment" escrito em microserviço (módulo separado)
* Comunicação assíncrona por mensageria
* Uma Saga para coordenar um fluxo distribuído
* Logs estruturados que permitam acompanhar o fluxo entre os serviços.

Endpoints implementados atualmente:

### Products

```
GET  /products
GET  /products/{id}
POST /products
```

### Orders

```
GET  /orders
GET  /pedidos?limit=10&offset=0
GET  /orders/{id}
POST /orders
POST /orders/{id}/pay
POST /orders/{id}/cancel
```
### Customers

```
GET  /customers
GET  /customers/{id}
POST /customers
```

## Estrutura do Serviço

```text

├── cmd
│   ├── app
│   │   └── main.go
│   └── seed
│       └── main.go
├── config
│   └── config.go
├── coverage.out
├── go.mod
├── go.sum
├── internal
│   ├── controllers
│   │   ├── customer_controller.go
│   │   ├── customer_controller_test.go
│   │   ├── errors.go
│   │   ├── errors_test.go
│   │   ├── mocks_test.go
│   │   ├── order_controller.go
│   │   ├── order_controller_test.go
│   │   ├── product_controller.go
│   │   └── product_controller_test.go
│   ├── database
│   │   └── postgres.go
│   ├── domain
│   │   ├── customer.go
│   │   ├── customer_test.go
│   │   ├── errors.go
│   │   ├── order.go
│   │   ├── order_test.go
│   │   ├── product.go
│   │   └── product_test.go
│   ├── dto
│   │   ├── customer_dto.go
│   │   ├── order_dto.go
│   │   └── product_dto.go
│   ├── events
│   │   └── events.go
│   ├── logger
│   │   ├── logger.go
│   │   └── slog_logger.go
│   ├── messaging
│   │   ├── handler.go
│   │   └── rabbitmq.go
│   ├── pagination
│   │   └── pagination.go
│   ├── repository
│   │   ├── customer_id_generator.go
│   │   ├── customer_repository.go
│   │   ├── dbtx.go
│   │   ├── logging_customer_repository.go
│   │   ├── logging_order_repository.go
│   │   ├── logging_order_repository_test.go
│   │   ├── logging_product_repository.go
│   │   ├── logging_product_repository_test.go
│   │   ├── mocks_test.go
│   │   ├── order_id_generator.go
│   │   ├── order_repository.go
│   │   ├── postgres_customer_id_generator.go
│   │   ├── postgres_customer_repository.go
│   │   ├── postgres_order_id_generator.go
│   │   ├── postgres_order_repository.go
│   │   ├── postgres_product_repository.go
│   │   ├── product_repository.go
│   │   └── unit_of_work.go
│   ├── routes
│   │   ├── customer_routes.go
│   │   ├── order_routes.go
│   │   ├── product_routes.go
│   │   └── routes.go
│   └── service
│       ├── customer_service.go
│       ├── customer_service_test.go
│       ├── mocks_test.go
│       ├── order_filters.go
│       ├── order_service.go
│       ├── order_service_test.go
│       ├── product_service.go
│       └── product_service_test.go
├── Makefile
├── migrations
│   ├── 000001_create_products.down.sql
│   ├── 000001_create_products.up.sql
│   ├── 000002_create_orders.down.sql
│   ├── 000002_create_orders.up.sql
│   ├── 000003_create_order_items.down.sql
│   ├── 000003_create_order_items.up.sql
│   ├── 000004_create_orders_id_seq.down.sql
│   ├── 000004_create_orders_id_seq.up.sql
│   ├── 000005_create_customers.down.sql
│   ├── 000005_create_customers.up.sql
│   ├── 000006_create_customers_id_seq.down.sql
│   ├── 000006_create_customers_id_seq.up.sql
│   ├── 000007_add_customer_id_to_orders.down.sql
│   └── 000007_add_customer_id_to_orders.up.sql
├── README.md
└── seeds
    ├── 100001_seed_products.down.sql
    └── 100001_seed_products.up.sql
```

## Executar a aplicação

### Rodando


Na raiz do projeto execute:

```bash
make run
```

### Migrations

```bash
make migrate-up
```

A API será iniciada localmente e ficará disponível para consumo através dos endpoints implementados.

## Funcionalidades demonstradas

Durante a evolução do projeto são demonstrados:

* Criação e gerenciamento de produtos
* Consulta de produtos por ID
* Consulta de pedidos por ID
* Listagem de produtos
* Listagem de pedidos
* Atualização de estoque
* Pagamento de pedidos
* Cancelamento com restauração de estoque
* Listagem paginada de pedidos
* Persistência dos dados em PostgreSQL
* Integração da aplicação com banco de dados utilizando Docker
* Criação de pedidos com transação
* Cadastrar clientes
* Listar clientes
* Criar pedido para um cliente
* Buscar cliente por id
* Pelo penos uma capacidade do sistema será escrita para um microserviço
* Comunicação assíncrona por mensageria
* Uma Saga para coordenar um fluxo distribuído
* Logs estruturados que permitam acompanhar o fluxo entre os serviços.

