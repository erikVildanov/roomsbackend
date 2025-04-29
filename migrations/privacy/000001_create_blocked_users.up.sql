CREATE TABLE blocked_users (
    blocker_id UUID NOT NULL,
    blocked_id UUID NOT NULL,
    blocked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blocker_id, blocked_id)
);