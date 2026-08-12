// Package access contains internal B3 token primitives. It registers no HTTP routes.
package access

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	TokenBytes           = 32
	VerificationTokenTTL = 15 * time.Minute
	StatusAccessPurpose  = "status_access"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, TokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}
