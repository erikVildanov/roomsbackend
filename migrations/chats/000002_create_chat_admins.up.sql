CREATE TABLE chat_admins (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);