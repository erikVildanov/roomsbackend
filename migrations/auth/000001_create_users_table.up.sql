CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    session_ttl TEXT DEFAULT '168h',
    created_at TIMESTAMP DEFAULT now()
);