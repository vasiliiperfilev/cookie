CREATE TABLE IF NOT EXISTS user_type (
    user_type_id smallserial PRIMARY KEY,  
    type_name text NOT NULL
);

ALTER TABLE user_type ADD CONSTRAINT user_type_type_name_max_len_check CHECK (char_length(type_name) <= 255);
ALTER TABLE user_type ADD CONSTRAINT user_type_type_name_min_len_check CHECK (char_length(type_name) > 0);