CREATE TABLE messages (
    id UUID PRIMARY KEY,
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    content TEXT,
    reply_to_id UUID REFERENCES messages(id),
    thread_chat_id UUID REFERENCES chats(id),
    pinned BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,
    edited_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);