CREATE TABLE starred_messages (
    user_id UUID NOT NULL,
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    starred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, message_id)
);