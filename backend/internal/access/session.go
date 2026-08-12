package access

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	BrowserSessionTTL        = 7 * 24 * time.Hour
	MaxActiveBrowserSessions = 5
)

var ErrSessionUnavailable = errors.New("status session unavailable")

// TransactionDB is satisfied by pgx pools. Session creation serializes per
// entry so the active-session limit remains correct under concurrent exchanges.
type TransactionDB interface {
	DB
	Begin(context.Context) (pgx.Tx, error)
}

type BrowserSession struct {
	EntryID   string
	RawID     string
	ExpiresAt time.Time
}

// CreateBrowserSession creates a fixed-lifetime session. The raw session ID is
// returned only to a trusted caller; PostgreSQL receives its SHA-256 hash.
func CreateBrowserSession(ctx context.Context, db TransactionDB, entryID string, now time.Time) (BrowserSession, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return BrowserSession{}, err
	}
	defer tx.Rollback(ctx)
	session, err := createBrowserSession(ctx, tx, entryID, now)
	if err != nil {
		return BrowserSession{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return BrowserSession{}, err
	}
	return session, nil
}

func createBrowserSession(ctx context.Context, tx pgx.Tx, entryID string, now time.Time) (BrowserSession, error) {
	// Locking the entry serializes the max-five policy for this entry.
	var lockedEntryID string
	if err := tx.QueryRow(ctx, `SELECT id::text FROM waitlist_entries WHERE id = $1 FOR UPDATE`, entryID).Scan(&lockedEntryID); err != nil {
		return BrowserSession{}, err
	}
	// When creating the sixth active session, revoke the oldest first. The
	// entry lock above makes this deterministic even if exchanges arrive at the
	// same time. The query also repairs any pre-existing excess in one pass.
	if _, err := tx.Exec(ctx, `
		WITH active_sessions AS (
			SELECT
				id,
				row_number() OVER (ORDER BY created_at ASC, id ASC) AS ordinal,
				count(*) OVER () AS total
			FROM status_access_sessions
			WHERE entry_id = $1 AND revoked_at IS NULL AND expires_at > $2
		)
		UPDATE status_access_sessions
		SET revoked_at = $2
		WHERE id IN (
			SELECT id
			FROM active_sessions
			WHERE total >= $3 AND ordinal <= total - ($3 - 1)
		)
	`, entryID, now, MaxActiveBrowserSessions); err != nil {
		return BrowserSession{}, err
	}
	raw, err := GenerateToken()
	if err != nil {
		return BrowserSession{}, err
	}
	expiresAt := now.Add(BrowserSessionTTL)
	if _, err := tx.Exec(ctx, `INSERT INTO status_access_sessions (entry_id, session_hash, expires_at) VALUES ($1, $2, $3)`, entryID, HashToken(raw), expiresAt); err != nil {
		return BrowserSession{}, err
	}
	return BrowserSession{EntryID: lockedEntryID, RawID: raw, ExpiresAt: expiresAt}, nil
}

// LookupBrowserSession returns only the internal entry identity for a valid,
// unrevoked session. It deliberately does not update expiry or last-used time.
func LookupBrowserSession(ctx context.Context, db DB, raw string, now time.Time) (string, error) {
	var entryID string
	err := db.QueryRow(ctx, `SELECT entry_id::text FROM status_access_sessions WHERE session_hash = $1 AND revoked_at IS NULL AND expires_at > $2`, HashToken(raw), now).Scan(&entryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrSessionUnavailable
	}
	return entryID, err
}

func RevokeBrowserSession(ctx context.Context, db DB, raw string, now time.Time) error {
	_, err := db.Exec(ctx, `UPDATE status_access_sessions SET revoked_at = COALESCE(revoked_at, $2) WHERE session_hash = $1`, HashToken(raw), now)
	return err
}

func RevokeAllBrowserSessions(ctx context.Context, db DB, entryID string, now time.Time) error {
	_, err := db.Exec(ctx, `UPDATE status_access_sessions SET revoked_at = COALESCE(revoked_at, $2) WHERE entry_id = $1`, entryID, now)
	return err
}

func CleanupExpiredBrowserSessions(ctx context.Context, db DB, now time.Time) error {
	_, err := db.Exec(ctx, `DELETE FROM status_access_sessions WHERE expires_at <= $1`, now)
	return err
}

// ConsumeVerificationAndCreateBrowserSession is intentionally internal. The
// verification token update and session insert share one transaction, so a
// failed session creation never consumes a usable verification token.
func ConsumeVerificationAndCreateBrowserSession(ctx context.Context, db TransactionDB, rawVerificationToken string, now time.Time) (BrowserSession, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return BrowserSession{}, err
	}
	defer tx.Rollback(ctx)
	entryID, err := ConsumeVerificationToken(ctx, tx, rawVerificationToken, now)
	if err != nil {
		return BrowserSession{}, err
	}
	session, err := createBrowserSession(ctx, tx, entryID, now)
	if err != nil {
		return BrowserSession{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return BrowserSession{}, err
	}
	return session, nil
}
