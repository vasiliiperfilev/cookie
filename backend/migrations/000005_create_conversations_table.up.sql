CREATE TABLE IF NOT EXISTS conversations (
    conversation_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    last_message_id bigint DEFAULT 0 REFERENCES messages(message_id) ON DELETE CASCADE,
    version integer NOT NULL DEFAULT 1
);