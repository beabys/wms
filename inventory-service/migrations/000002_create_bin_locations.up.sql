CREATE TABLE IF NOT EXISTS bin_locations (
    id TEXT PRIMARY KEY,
    warehouse_zone TEXT NOT NULL,
    aisle TEXT NOT NULL,
    rack TEXT NOT NULL,
    shelf TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(warehouse_zone, aisle, rack, shelf)
);
