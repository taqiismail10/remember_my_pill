package access

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB is satisfied by pgx pools and transactions without introducing an ORM.
type DB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type VerificationToken struct {
	ID        string
	EntryID   string
	ExpiresAt time.Time
}

// CreateVerificationToken revokes prior active tokens for this entry/purpose
// before recording only the new token hash. The raw token is never persisted.
func CreateVerificationToken(ctx context.Context, db DB, entryID, raw string, now time.Time) (VerificationToken, error) {
	if _, err := db.Exec(ctx, `UPDATE waitlist_verification_tokens SET revoked_at = $1 WHERE entry_id = $2 AND purpose = $3 AND used_at IS NULL AND revoked_at IS NULL`, now, entryID, StatusAccessPurpose); err != nil {
		return VerificationToken{}, err
	}
	var token VerificationToken
	err := db.QueryRow(ctx, `INSERT INTO waitlist_verification_tokens (entry_id, token_hash, purpose, expires_at) VALUES ($1, $2, $3, $4) RETURNING id::text, entry_id::text, expires_at`, entryID, HashToken(raw), StatusAccessPurpose, now.Add(VerificationTokenTTL)).Scan(&token.ID, &token.EntryID, &token.ExpiresAt)
	return token, err
}

// ConsumeVerificationToken is atomic: an expired, used, or revoked token has
// no returned row and cannot be consumed again.
func ConsumeVerificationToken(ctx context.Context, db DB, raw string, now time.Time) (string, error) {
	var entryID string
	err := db.QueryRow(ctx, `UPDATE waitlist_verification_tokens SET used_at = $1 WHERE token_hash = $2 AND purpose = $3 AND expires_at > $1 AND used_at IS NULL AND revoked_at IS NULL RETURNING entry_id::text`, now, HashToken(raw), StatusAccessPurpose).Scan(&entryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrTokenUnavailable
	}
	return entryID, err
}

var ErrTokenUnavailable = errors.New("verification token unavailable")

// RotateStatusToken replaces the stored hash and returns raw material only to
// its trusted internal caller.
func RotateStatusToken(ctx context.Context, db DB, entryID string) (string, error) {
	raw, err := GenerateToken()
	if err != nil {
		return "", err
	}
	result, err := db.Exec(ctx, `UPDATE waitlist_entries SET status_token_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, HashToken(raw), entryID)
	if err != nil {
		return "", err
	}
	if result.RowsAffected() != 1 {
		return "", pgx.ErrNoRows
	}
	return raw, nil
}
