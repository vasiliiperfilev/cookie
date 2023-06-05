CREATE TABLE IF NOT EXISTS users (
    user_id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    user_type_id smallint NOT NULL,
    image_id varchar(255) NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE users ADD CONSTRAINT fk_users_user_types FOREIGN KEY (user_type_id) REFERENCES user_types (user_type_id);

INSERT INTO users (user_id, email, password_hash, user_type_id, image_id)
VALUES 
    (0, 'sentinel@user', 'hash', 1, 'sentinel');