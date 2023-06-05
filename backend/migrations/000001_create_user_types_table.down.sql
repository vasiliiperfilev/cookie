ALTER TABLE user_types DROP CONSTRAINT IF EXISTS user_types_type_name_max_len_check;
ALTER TABLE user_types DROP CONSTRAINT IF EXISTS user_types_type_name_min_len_check;

DROP TABLE IF EXISTS user_types;