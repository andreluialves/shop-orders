CREATE SEQUENCE IF NOT EXISTS customers_id_seq;

-- Ajusta a sequência para continuar depois do maior ID já existente,
-- evitando colisão com clientes já cadastrados (ex: CUST-001).
-- "CUST-" tem 5 caracteres, então a parte numérica começa na posição 6.
SELECT setval(
	'customers_id_seq',
	COALESCE(
		(SELECT MAX(CAST(SUBSTRING(id FROM 6) AS INTEGER)) FROM customers),
		0
	)
);