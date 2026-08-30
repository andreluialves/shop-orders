# Shop Orders — Mono-repo

Projeto desenvolvido em Go para estudo de arquitetura em camadas, boas práticas de desenvolvimento e criação de APIs REST. O sistema de pedidos composto por dois serviços independentes, comunicando-se
de forma assíncrona via RabbitMQ (padrão Saga, coreografia).

## Tecnologias

* Go
* PostgreSQL
* Docker
* REST API

## Serviços

| Serviço | Descrição | Documentação |
|---|---|---|
| [`shop-orders`](./shop-orders) | Pedidos, produtos e clientes | [README](./shop-orders/README.md) |
| [`payments-service`](./payments-service) | Processamento de pagamentos | [README](./payments-service/README.md) |

## Arquitetura

shop-orders --publica--> payment.requested --consome--> payments-service
shop-orders <--consome-- payment.processed <--publica-- payments-service


Cada serviço:
- É compilado e implantado de forma independente
- Possui seu próprio banco de dados e migrations
- Não acessa diretamente os dados do outro serviço
- Se comunica exclusivamente via eventos assíncronos (RabbitMQ)

## Estrutura dp Projeto

```text

├── payments-service // Pasta do serviço de pedidos, produtos e clientes
├── shop-orders      // Pasta do serviço de processamento de pagamentos
├── README.md
└── docker-compose.yml
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

## Rodando localmente

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
### Infraestrutura compartilhada

Sobe a infraestrutura compartilhada (Postgres de cada serviço + RabbitMQ):

```bash
docker-compose up -d
```

Depois, em terminais separados:

```bash
cd shop-orders && make run
cd payments-service && make run
```

>**Consulte o README de cada serviço para detalhes de configuração e migrations.**

## Acompanhamento do Projeto

O desenvolvimento do projeto é acompanhado através de um **Project Board** do GitHub, onde são organizadas as tarefas, melhorias, correções e novas funcionalidades planejadas.

O board pode ser acessado através do link:

[GitHub Project Board](https://github.com/users/andreluialves/projects/9/views/2)