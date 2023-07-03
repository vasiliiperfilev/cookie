CREATE TABLE IF NOT EXISTS orders_states(
    order_id bigint NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    state_id bigint NOT NULL REFERENCES states(state_id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY(order_id, state_id)
);
