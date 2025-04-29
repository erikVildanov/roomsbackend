CREATE TABLE chats (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL, -- 'direct', 'group', 'channel'
    name TEXT,
    created_by UUID NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    invite_token TEXT,
    invite_expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);