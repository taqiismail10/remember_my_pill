ALTER TABLE waitlist_entries
    ADD COLUMN marketing_consented_at TIMESTAMPTZ NULL,
    ADD COLUMN marketing_consent_version VARCHAR(30) NULL,
    ADD COLUMN marketing_withdrawn_at TIMESTAMPTZ NULL,
    ADD CONSTRAINT waitlist_entries_marketing_consent_evidence CHECK (
        (marketing_consented_at IS NULL AND marketing_consent_version IS NULL)
        OR (marketing_consented_at IS NOT NULL AND marketing_consent_version IS NOT NULL)
    );
