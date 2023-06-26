CREATE TABLE IF NOT EXISTS units (
    unit_id serial PRIMARY KEY,  
    name VARCHAR(255) NOT NULL
);

ALTER TABLE units ADD CONSTRAINT units_name_min_length CHECK (char_length(name) > 0);

INSERT INTO units (unit_id, name)
VALUES 
    (1, 'l'),
    (2, 'kg');