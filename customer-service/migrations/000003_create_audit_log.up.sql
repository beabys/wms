CREATE TABLE IF NOT EXISTS customer_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    action VARCHAR(50) NOT NULL,
    performed_by VARCHAR(255) NOT NULL,
    details TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_customer ON customer_audit_log(customer_id);
CREATE INDEX idx_audit_log_created ON customer_audit_log(created_at DESC);
