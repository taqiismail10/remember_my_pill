CREATE TABLE referral_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  referrer_id UUID NOT NULL REFERENCES waitlist_entries(id),
  referred_entry_id UUID NOT NULL REFERENCES waitlist_entries(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT referral_events_referred_entry_unique UNIQUE (referred_entry_id),
  CONSTRAINT referral_events_not_self_referred CHECK (referrer_id <> referred_entry_id)
);

CREATE INDEX referral_events_referrer_id_idx ON referral_events (referrer_id);

CREATE TABLE waitlist_verification_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entry_id UUID NOT NULL REFERENCES waitlist_entries(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  purpose TEXT NOT NULL CHECK (purpose = 'status_access'),
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX waitlist_verification_tokens_expiry_idx ON waitlist_verification_tokens (expires_at);
CREATE INDEX waitlist_verification_tokens_active_lookup_idx ON waitlist_verification_tokens (token_hash) WHERE used_at IS NULL AND revoked_at IS NULL;
CREATE INDEX waitlist_verification_tokens_entry_recent_idx ON waitlist_verification_tokens (entry_id, created_at DESC);
