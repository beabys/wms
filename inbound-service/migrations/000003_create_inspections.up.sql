CREATE TABLE IF NOT EXISTS inspections (
    id            TEXT PRIMARY KEY,
    inbound_id    TEXT NOT NULL UNIQUE REFERENCES inbounds(id) ON DELETE CASCADE,
    inspector_id  TEXT NOT NULL,
    notes         TEXT NOT NULL DEFAULT '',
    photos_json   TEXT NOT NULL DEFAULT '[]',
    passed        BOOLEAN NOT NULL DEFAULT FALSE,
    inspected_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_inspections_inbound_id ON inspections(inbound_id);
