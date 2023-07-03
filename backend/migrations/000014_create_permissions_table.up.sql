CREATE TABLE IF NOT EXISTS permissions(
    permission_id smallserial PRIMARY KEY,
    name varchar(255) NOT NULL
);

INSERT INTO permissions (name)
VALUES 
    ('order:create'),
    ('order:accept'),
    ('order:decline'),
    ('order:fulfill'),
    ('order:confirm_fulfillment'),
    ('order:supplier_changes'),
    ('order:client_changes');