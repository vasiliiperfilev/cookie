CREATE TABLE IF NOT EXISTS app_user (
    user_id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    user_type_id smallint NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE app_user ADD CONSTRAINT fk_app_user_user_type FOREIGN KEY (user_type_id) REFERENCES user_type (user_type_id);