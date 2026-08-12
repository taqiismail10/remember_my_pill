ALTER TABLE waitlist_entries
  DROP CONSTRAINT IF EXISTS waitlist_entries_not_self_referred,
  DROP CONSTRAINT IF EXISTS waitlist_entries_status_token_hash_unique,
  DROP CONSTRAINT IF EXISTS waitlist_entries_referral_code_unique,
  DROP COLUMN IF EXISTS updated_at,
  DROP COLUMN IF EXISTS consented_at,
  DROP COLUMN IF EXISTS consent_version,
  DROP COLUMN IF EXISTS status_token_hash,
  DROP COLUMN IF EXISTS referred_by_id,
  DROP COLUMN IF EXISTS referral_code;
