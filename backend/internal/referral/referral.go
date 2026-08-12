// Package referral contains internal B3 primitives; it does not alter signup behavior.
package referral

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const CodeBytes = 16

func GenerateCode() (string, error) {
	bytes := make([]byte, CodeBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate referral code: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateUnique retries a bounded number of random codes; the database
// unique constraint remains the final concurrency safeguard.
func GenerateUnique(isAvailable func(string) (bool, error)) (string, error) {
	for attempt := 0; attempt < 5; attempt++ {
		code, err := GenerateCode()
		if err != nil {
			return "", err
		}
		available, err := isAvailable(code)
		if err != nil {
			return "", err
		}
		if available {
			return code, nil
		}
	}
	return "", ErrCollisionLimit
}

var ErrCollisionLimit = fmt.Errorf("referral code collision limit reached")

type EventStore interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

var ErrAttributionUnavailable = errors.New("referral attribution unavailable")

// CreateEvent is the internal, database-constrained primitive for later
// activation. It is not called by the B3A production signup path.
func CreateEvent(ctx context.Context, db EventStore, referrerID, referredEntryID string) error {
	if referrerID == "" || referredEntryID == "" || referrerID == referredEntryID {
		return ErrAttributionUnavailable
	}
	_, err := db.Exec(ctx, `INSERT INTO referral_events (referrer_id, referred_entry_id) VALUES ($1, $2)`, referrerID, referredEntryID)
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23514") {
		return ErrAttributionUnavailable
	}
	return err
}
