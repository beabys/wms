CREATE TABLE IF NOT EXISTS customers (
    id            TEXT PRIMARY KEY,
    company_name  TEXT NOT NULL,
    vat_number    TEXT NOT NULL UNIQUE,
    address       JSONB NOT NULL DEFAULT '{}',
    status        TEXT NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending', 'active', 'suspended')),
    rate_card_id  TEXT NOT NULL DEFAULT '',
    credit_balance BIGINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_status ON customers (status);
CREATE INDEX idx_customers_vat_number ON customers (vat_number);
