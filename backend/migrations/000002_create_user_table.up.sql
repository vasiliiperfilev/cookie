CREATE TABLE IF NOT EXISTS users (
    user_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    user_type_id smallint NOT NULL,
    image_id varchar(255) NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE users ADD CONSTRAINT fk_users_user_types FOREIGN KEY (user_type_id) REFERENCES user_types (user_type_id);