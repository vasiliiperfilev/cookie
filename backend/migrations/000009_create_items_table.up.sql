CREATE TABLE IF NOT EXISTS items (
    item_id bigserial PRIMARY KEY,
    supplier_id bigint REFERENCES users(user_id) ON DELETE CASCADE,
    unit_id bigint REFERENCES units(unit_id),
    name varchar(255) NOT NULL,
    image_url varchar(255) NOT NULL
);