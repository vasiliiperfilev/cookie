CREATE TABLE IF NOT EXISTS orders (
    order_id bigserial PRIMARY KEY,
    messsage_id bigint REFERENCES messages(message_id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    supplier_comment text NOT NULL DEFAULT '',
    client_comment text NOT NULL DEFAULT ''
);