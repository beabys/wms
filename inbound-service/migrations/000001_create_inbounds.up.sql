CREATE TABLE IF NOT EXISTS inbounds (
    id               TEXT PRIMARY KEY,
    customer_id      TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'submitted',
    expected_date    TEXT NOT NULL,
    packaging_profile_id TEXT,
    notes            TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inbounds_customer_id ON inbounds(customer_id);
CREATE INDEX idx_inbounds_status ON inbounds(status);
CREATE INDEX idx_inbounds_expected_date ON inbounds(expected_date);
