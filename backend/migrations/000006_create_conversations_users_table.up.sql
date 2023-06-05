CREATE TABLE IF NOT EXISTS conversations_users (
    conversation_id bigint NOT NULL REFERENCES conversations(conversation_id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    last_read_message_id bigint DEFAULT 0 REFERENCES messages(message_id) ON DELETE CASCADE,
    PRIMARY KEY(conversation_id, user_id)
);

