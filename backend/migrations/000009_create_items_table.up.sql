CREATE TABLE IF NOT EXISTS items (
    item_id bigserial PRIMARY KEY,
    supplier_id bigint REFERENCES users(user_id) ON DELETE CASCADE,
    unit_id bigint REFERENCES units(unit_id),
    size numeric(7,2) NOT NULL,
    name varchar(255) NOT NULL,
    image_url varchar(255) NOT NULL
);

ALTER TABLE items ADD CONSTRAINT items_name_min_length CHECK (char_length(name) > 0);
ALTER TABLE items ADD CONSTRAINT items_size_positive CHECK (size > 0);