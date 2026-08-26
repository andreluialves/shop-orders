CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

-- Ajusta a sequência para começar depois do maior ID já existente no banco,
-- evitando colisão com pedidos criados antes dessa migration (ex: PED-001).
-- "PED-" tem 4 caracteres, então a parte numérica começa na posição 5.
SELECT setval(
	'orders_id_seq',
	COALESCE(
		(SELECT MAX(CAST(SUBSTRING(id FROM 5) AS INTEGER)) FROM orders),
		0
	)
);