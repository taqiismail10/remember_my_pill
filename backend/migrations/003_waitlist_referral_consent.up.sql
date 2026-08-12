ALTER TABLE waitlist_entries
  ADD COLUMN referral_code VARCHAR(50),
  ADD COLUMN referred_by_id UUID REFERENCES waitlist_entries(id),
  ADD COLUMN status_token_hash TEXT,
  ADD COLUMN consent_version VARCHAR(30),
  ADD COLUMN consented_at TIMESTAMPTZ,
  ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ADD CONSTRAINT waitlist_entries_referral_code_unique UNIQUE (referral_code),
  ADD CONSTRAINT waitlist_entries_status_token_hash_unique UNIQUE (status_token_hash),
  ADD CONSTRAINT waitlist_entries_not_self_referred CHECK (referred_by_id IS NULL OR referred_by_id <> id);
