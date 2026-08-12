CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE waitlist_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 1 AND 100),
  email_normalized VARCHAR(255) NOT NULL UNIQUE CHECK (char_length(email_normalized) BETWEEN 3 AND 255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
