# Shop Orders — Mono-repo

Sistema de pedidos composto por dois serviços independentes, comunicando-se
de forma assíncrona via RabbitMQ (padrão Saga, coreografia).

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

## Rodando localmente

Sobe a infraestrutura compartilhada (Postgres de cada serviço + RabbitMQ):

```bash
docker-compose up -d
```

Depois, em terminais separados:

```bash
cd shop-orders && make run
cd payments-service && make run
```

Consulte o README de cada serviço para detalhes de configuração e migrations.