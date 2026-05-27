CREATE TABLE IF NOT EXISTS inbound_items (
    id                TEXT PRIMARY KEY,
    inbound_id        TEXT NOT NULL REFERENCES inbounds(id) ON DELETE CASCADE,
    sku               TEXT NOT NULL,
    quantity_declared INTEGER NOT NULL,
    quantity_received INTEGER NOT NULL DEFAULT 0,
    dimensions        TEXT NOT NULL DEFAULT '',
    weight            DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE INDEX idx_inbound_items_inbound_id ON inbound_items(inbound_id);
