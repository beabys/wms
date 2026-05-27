CREATE TABLE IF NOT EXISTS stock_entries (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES products(id),
    bin_location_id TEXT REFERENCES bin_locations(id),
    quantity DOUBLE PRECISION NOT NULL DEFAULT 0,
    reserved_quantity DOUBLE PRECISION NOT NULL DEFAULT 0,
    lot_number TEXT NOT NULL DEFAULT '',
    expiry_date TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'active',
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stock_entries_product_id ON stock_entries(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_entries_bin_location_id ON stock_entries(bin_location_id);
CREATE INDEX IF NOT EXISTS idx_stock_entries_status ON stock_entries(status);
