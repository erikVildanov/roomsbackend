CREATE TABLE files (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    filename TEXT NOT NULL,
    mime_type TEXT,
    size INT,
    url TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);