# Shop Orders

Projeto desenvolvido em Go para estudo de arquitetura em camadas, boas práticas de desenvolvimento e regras de negócio aplicadas ao gerenciamento de pedidos.

## Tecnologias

* Go

## Funcionalidades

* Cadastro de produtos
* Criação de pedidos
* Validação de cliente e itens
* Controle de estoque
* Pagamento de pedidos
* Cancelamento de pedidos com restauração do estoque
* Listagem de pedidos
* Filtro de pedidos por:

  * Status (pendentes e pagos)
  * Valor mínimo do pedido

## Estrutura do Projeto

```text
cmd
├── app/
|   ├── main.go
internal/
├── domain/
│   ├── order.go
│   ├── product.go
│   └── errors.go
│
├── repository/
│   ├── product_repository.go
│   ├── order_repository.go
│   ├── memory_product_repository.go
│   └── memory_order_repository.go
│
└── service/
    ├── order_service.go
    └── order_filters.go
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
* Git

Verifique a versão do Go instalada:

```bash
go version
```

### Instalar dependências

Caso existam dependências externas, execute:

```bash
go mod tidy
```

### Executar a aplicação

Na raiz do projeto execute:

```bash
go run ./cmd/app/
```

ou

```bash
go run ./cmd/app/main.go
```

A aplicação será executada no terminal, demonstrando o fluxo completo de criação e gerenciamento de pedidos.

## Funcionalidades demonstradas

O `main.go` apresenta exemplos de:

* Cadastro de produtos
* Criação de pedidos
* Validação de estoque
* Atualização de estoque após a compra
* Pagamento de pedidos
* Cancelamento de pedidos
* Busca de pedidos por ID
* Listagem de pedidos
* Filtros por status e valor

## Próxima etapa (Fase 2)

A segunda fase do projeto terá como objetivo substituir os repositórios em memória por persistência utilizando PostgreSQL e Docker.
