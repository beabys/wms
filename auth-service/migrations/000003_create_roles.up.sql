CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed system roles
INSERT INTO roles (name, description, is_system) VALUES
('admin', 'Full system access', true),
('warehouse_staff', 'Inbound, inventory, orders, shipping', true),
('billing_manager', 'Billing, ledger, reports', true),
('customer', 'Own data only (orders, inbounds, returns, ledger)', true)
ON CONFLICT (name) DO NOTHING;
