CREATE TABLE IF NOT EXISTS states (
    state_id smallserial PRIMARY KEY,  
    name varchar(255) NOT NULL
);

INSERT INTO states (name)
VALUES 
    ('created'),
    ('accepted'),
    ('declined'),
    ('fulfilled'),
    ('confirmed_fulfillment'),
    ('supplier_changes'),
    ('client_changes');