CREATE TABLE IF NOT EXISTS orders_items (
    order_id bigint NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    item_id bigint NOT NULL REFERENCES items(item_id) ON DELETE CASCADE,
    quantity int NOT NULL DEFAULT 1,
    PRIMARY KEY(order_id, item_id)
);
