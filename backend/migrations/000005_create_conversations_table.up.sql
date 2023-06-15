CREATE TABLE IF NOT EXISTS conversations (
    conversation_id bigserial PRIMARY KEY,
    last_message_id bigint DEFAULT 0 REFERENCES messages(message_id) ON DELETE CASCADE,
    version integer NOT NULL DEFAULT 1
);

INSERT INTO conversations (conversation_id)
VALUES 
    (0);

ALTER TABLE messages ADD CONSTRAINT fk_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE;