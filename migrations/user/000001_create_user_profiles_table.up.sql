CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    bio TEXT,
    position TEXT,
    phone_number TEXT,
    created_at TIMESTAMP DEFAULT now()
);
