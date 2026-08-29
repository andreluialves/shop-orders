-- Adiciona a nova coluna, ainda sem NOT NULL, pra permitir o backfill
ALTER TABLE orders ADD COLUMN customer_id VARCHAR(20) REFERENCES customers(id);

-- Melhor esforço: cada pedidos existentes com clientes pelo nome exato.
-- Pedidos cujo nome não bate com nenhum cliente cadastrado ficam com
-- customer_id NULL, e precisam ser corrigidos manualmente se necessário.
UPDATE orders o
SET customer_id = c.id
FROM customers c
WHERE c.name = o.customer;

-- Remove a coluna antiga, que guardava o nome digitado livremente
ALTER TABLE orders DROP COLUMN customer;