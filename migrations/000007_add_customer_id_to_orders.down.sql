ALTER TABLE orders ADD COLUMN customer VARCHAR(255);

UPDATE orders o
SET customer = c.name
FROM customers c
WHERE c.id = o.customer_id;

ALTER TABLE orders DROP COLUMN customer_id;