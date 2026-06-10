ALTER TABLE invite_tokens ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';
UPDATE invite_tokens SET status = 'used' WHERE used = true;
UPDATE invite_tokens SET status = 'cancelled' WHERE used = false AND expires_at < NOW();
-- Deduplicate: for emails with multiple pending invites, cancel all but the most recent
UPDATE invite_tokens t SET status = 'cancelled'
  FROM (
    SELECT email, created_at, ROW_NUMBER() OVER (PARTITION BY email ORDER BY created_at DESC) AS rn
    FROM invite_tokens WHERE status = 'pending'
  ) dups
  WHERE t.email = dups.email AND t.created_at = dups.created_at AND dups.rn > 1 AND t.status = 'pending';
ALTER TABLE invite_tokens DROP COLUMN used;
CREATE UNIQUE INDEX idx_invite_tokens_pending_email ON invite_tokens(email) WHERE status = 'pending';
