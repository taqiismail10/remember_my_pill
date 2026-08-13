package access

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// StatusAccessEligibility is deliberately internal. Callers must keep every
// outcome externally indistinguishable to avoid email enumeration.
type StatusAccessEligibility string

const (
	StatusAccessEligible   StatusAccessEligibility = "eligible"
	StatusAccessIneligible StatusAccessEligibility = "ineligible"
)

// EvaluateStatusAccessEligibility permits only entries with explicit required
// waitlist consent. Legacy rows with NULL consent are never eligible.
func EvaluateStatusAccessEligibility(ctx context.Context, db DB, normalizedEmail string) (StatusAccessEligibility, string, error) {
	var entryID string
	err := db.QueryRow(ctx, `
		SELECT id::text
		FROM waitlist_entries
		WHERE email_normalized = $1
			AND consent_version IS NOT NULL
			AND consented_at IS NOT NULL
	`, normalizedEmail).Scan(&entryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return StatusAccessIneligible, "", nil
	}
	if err != nil {
		return StatusAccessIneligible, "", err
	}
	return StatusAccessEligible, entryID, nil
}

type WaitlistStatus struct {
	Rank          int64
	ReferralCount int64
	ReferralCode  string
}

// LookupWaitlistStatus returns only the approved non-identifying status
// fields. Rank ordering is stable on created_at followed by id.
func LookupWaitlistStatus(ctx context.Context, db DB, entryID string) (WaitlistStatus, error) {
	var status WaitlistStatus
	err := db.QueryRow(ctx, `
		SELECT
			1 + (
				SELECT count(*)
				FROM waitlist_entries earlier
				WHERE (earlier.created_at, earlier.id) < (entry.created_at, entry.id)
			),
			(SELECT count(*) FROM referral_events WHERE referrer_id = entry.id),
			COALESCE(entry.referral_code, '')
		FROM waitlist_entries entry
		WHERE entry.id = $1
	`, entryID).Scan(&status.Rank, &status.ReferralCount, &status.ReferralCode)
	return status, err
}
