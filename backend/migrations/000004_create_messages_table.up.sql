CREATE TABLE IF NOT EXISTS messages (
    message_id bigserial PRIMARY KEY,
    sender_id bigint REFERENCES users(user_id) ON DELETE CASCADE,
    conversation_id bigint NOT NULL, -- constrait is added when conversations relation created
    prev_message_id bigint DEFAULT 0 REFERENCES messages(message_id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    content citext NOT NULL
);

INSERT INTO messages (message_id, sender_id, conversation_id, content)
VALUES 
    (0, 0, 0, 'sentinel node');