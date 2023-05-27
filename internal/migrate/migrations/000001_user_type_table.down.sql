ALTER TABLE user_type DROP CONSTRAINT IF EXISTS user_type_type_name_max_len_check;
ALTER TABLE user_type DROP CONSTRAINT IF EXISTS user_type_type_name_min_len_check;

DROP TABLE IF EXISTS user_type;