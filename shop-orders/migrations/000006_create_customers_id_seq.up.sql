CREATE SEQUENCE IF NOT EXISTS customers_id_seq;

-- Ajusta a sequência para continuar depois do maior ID já existente,
-- evitando colisão com clientes já cadastrados (ex: CUST-001).
-- "CUST-" tem 5 caracteres, então a parte numérica começa na posição 6.
DO $$
DECLARE
	max_id INTEGER;
BEGIN
	SELECT MAX(CAST(SUBSTRING(id FROM 6) AS INTEGER)) INTO max_id FROM customers;

	IF max_id IS NULL THEN
		PERFORM setval('customers_id_seq', 1, false);
	ELSE
		PERFORM setval('customers_id_seq', max_id);
	END IF;
END $$;