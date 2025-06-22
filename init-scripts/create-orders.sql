CREATE TABLE orders
(
    order_id     SERIAL PRIMARY KEY,
    customer_id  INT            NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status       VARCHAR(50) DEFAULT 'Pending'
);

INSERT INTO orders (customer_id, total_amount)
VALUES (1, 2.3),
       (1, 3.3),
       (2, 1.3);
