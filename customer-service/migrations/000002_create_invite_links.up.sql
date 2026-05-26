CREATE TABLE IF NOT EXISTS invite_links (
    id         TEXT PRIMARY KEY,
    email      TEXT NOT NULL,
    token      TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invite_links_token ON invite_links (token);
CREATE INDEX idx_invite_links_email ON invite_links (email);
