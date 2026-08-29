# Payments Service

Serviço responsável por processar pagamentos de pedidos, de forma assíncrona.

## Responsabilidade

- Recebe solicitações de pagamento via evento `payment.requested`
- Decide aprovação/recusa e persiste o registro do pagamento
- Publica o resultado via evento `payment.processed`

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