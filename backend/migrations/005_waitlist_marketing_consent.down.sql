ALTER TABLE waitlist_entries
    DROP CONSTRAINT waitlist_entries_marketing_consent_evidence,
    DROP COLUMN marketing_withdrawn_at,
    DROP COLUMN marketing_consent_version,
    DROP COLUMN marketing_consented_at;
