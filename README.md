# Shop Orders

Projeto desenvolvido em Go para estudo de arquitetura em camadas, boas práticas de desenvolvimento, criação de APIs REST e aplicação de regras de negócio no gerenciamento de produtos e pedidos.

## Tecnologias

* Go
* PostgreSQL
* Docker
* REST API

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

Endpoints implementados atualmente:

### Products

```
GET /products
GET /products/{id}
POST /products
```

### Orders

```
GET /orders
GET /orders/{id}
```

*(Novos endpoints serão adicionados conforme a evolução da Fase 2.)*

## Estrutura do Projeto

```text
.
├── cmd
│   └── app
│       └── main.go
│
├── internal
│   ├── database
│   │   └── postgres.go
│   │
│   ├── domain
│   │   ├── order.go
│   │   ├── product.go
│   │   └── errors.go
│   │
│   ├── dto
│   │   ├── order_dto.go
│   │   └── product_dto.go
│   │
│   ├── repository
│   │   ├── product_repository.go
│   │   ├── order_repository.go
│   │   ├── memory_product_repository.go
│   │   └── memory_order_repository.go
│   │
│   ├── routes
│   │   ├── routes.go
│   │   ├── order_routes.go
│   │   └── product_routes.go
│   │
│   └── service
│       ├── order_service.go
│       └── order_filters.go
│
├── migrations
│   ├── 000001_create_products.down.sql
│   ├── 000001_create_products.up.sql
│   ├── 000002_create_orders.down.sql
│   └── 000002_create_orders.up.sql
│
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Como clonar o projeto

Clone o repositório utilizando o Git:

```bash
git clone https://github.com/andreluialves/shop-orders.git
```

Acesse a pasta do projeto:

```bash
cd shop-orders
```

## Como executar localmente

### Pré-requisitos

* Go 1.24 ou superior instalado
* Docker e Docker Compose
* Git

Verifique a versão do Go instalada:

```bash
go version
```

Verifique a instalação do Docker:

```bash
docker --version
```

## Configuração do banco de dados

O projeto utiliza PostgreSQL executado através de Docker para persistência dos dados.

Suba os containers:

```bash
docker compose up -d
```

O banco ficará disponível para a aplicação conforme as configurações definidas no arquivo de ambiente do projeto.

## Instalar dependências

Execute:

```bash
go mod tidy
```

## Executar a aplicação

Na raiz do projeto execute:

```bash
go run ./cmd/app/
```

ou

```bash
go run ./cmd/app/main.go
```

A API será iniciada localmente e ficará disponível para consumo através dos endpoints implementados.

## Funcionalidades demonstradas

Durante a evolução do projeto são demonstrados:

* Criação e gerenciamento de produtos
* Consulta de produtos por ID
* Consulta de pedidos por ID
* Listagem de produtos
* Listagem de pedidos
* Persistência dos dados em PostgreSQL
* Integração da aplicação com banco de dados utilizando Docker

## Próximas etapas

A Fase 2 encontra-se em andamento e terá como próximos objetivos:

* Criação de pedidos com transação
* Atualização de estoque
* Pagamento de pedidos
* Cancelamento com restauração de estoque
* Listagem paginada de pedidos
* Cadastrar clientes
* Listar clientes
* Buscar cliente por id;
* Melhorias na documentação da API
