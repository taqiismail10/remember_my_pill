CREATE TABLE status_access_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entry_id UUID NOT NULL REFERENCES waitlist_entries(id) ON DELETE CASCADE,
  session_hash TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  last_used_at TIMESTAMPTZ NULL
);

CREATE INDEX status_access_sessions_entry_active_idx
  ON status_access_sessions (entry_id, created_at ASC, id ASC)
  WHERE revoked_at IS NULL;
CREATE INDEX status_access_sessions_expiry_idx ON status_access_sessions (expires_at);
