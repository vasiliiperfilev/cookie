CREATE TABLE IF NOT EXISTS user_types (
    user_type_id smallserial PRIMARY KEY,  
    type_name text NOT NULL
);

ALTER TABLE user_types ADD CONSTRAINT user_types_type_name_max_len_check CHECK (char_length(type_name) <= 255);
ALTER TABLE user_types ADD CONSTRAINT user_types_type_name_min_len_check CHECK (char_length(type_name) > 0);

INSERT INTO user_types (user_type_id, type_name)
VALUES 
    (1, 'supplier'),
    (2, 'client');