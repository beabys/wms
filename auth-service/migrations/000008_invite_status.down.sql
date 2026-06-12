ALTER TABLE invite_tokens ADD COLUMN used BOOLEAN NOT NULL DEFAULT false;
UPDATE invite_tokens SET used = true WHERE status = 'used';
DROP INDEX IF EXISTS idx_invite_tokens_pending_email;
ALTER TABLE invite_tokens DROP COLUMN status;
