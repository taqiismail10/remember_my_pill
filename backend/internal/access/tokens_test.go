package access

import (
	"context"
	"regexp"
	"testing"
)

func TestGenerateTokenAndHash(t *testing.T) {
	first, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !regexp.MustCompile(`^[A-Za-z0-9_-]{43}$`).MatchString(first) {
		t.Fatalf("unexpected token %q", first)
	}
	hash := HashToken(first)
	if hash == first || hash != HashToken(first) || !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(hash) {
		t.Fatalf("unexpected hash %q", hash)
	}
}

func TestFakeEmailSender(t *testing.T) {
	fake := &FakeEmailSender{Result: DeliveryAccepted}
	if got := fake.SendStatusAccessEmail(context.Background(), StatusAccessEmail{To: "person@example.test", VerificationURL: "https://rmp.example.test/waitlist/verify#v=test"}); got != DeliveryAccepted || len(fake.Sent) != 1 {
		t.Fatalf("result=%q sent=%d", got, len(fake.Sent))
	}
}
