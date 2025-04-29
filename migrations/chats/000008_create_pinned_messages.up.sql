CREATE TABLE pinned_messages (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    pinned_by UUID NOT NULL,
    pinned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, message_id)
);