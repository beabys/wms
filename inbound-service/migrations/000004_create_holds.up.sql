CREATE TABLE IF NOT EXISTS holds (
    id            TEXT PRIMARY KEY,
    inbound_id    TEXT NOT NULL UNIQUE REFERENCES inbounds(id) ON DELETE CASCADE,
    reason        TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at   TIMESTAMPTZ,
    released_by   TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_holds_inbound_id ON holds(inbound_id);
