CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

-- Ajusta a sequência para começar depois do maior ID já existente no banco,
-- evitando colisão com pedidos criados antes dessa migration (ex: PED-001).
-- "PED-" tem 4 caracteres, então a parte numérica começa na posição 5.
DO $$
DECLARE
	max_id INTEGER;
BEGIN
	SELECT MAX(CAST(SUBSTRING(id FROM 5) AS INTEGER)) INTO max_id FROM orders;

	IF max_id IS NULL THEN
		PERFORM setval('orders_id_seq', 1, false); -- próximo nextval() = 1
	ELSE
		PERFORM setval('orders_id_seq', max_id); -- próximo nextval() = max_id + 1
	END IF;
END $$;