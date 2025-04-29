CREATE TABLE chat_participants (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    is_muted BOOLEAN DEFAULT FALSE,
    is_favorite BOOLEAN DEFAULT FALSE,
    last_read_at TIMESTAMP,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);