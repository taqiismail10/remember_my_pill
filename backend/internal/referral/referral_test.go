package referral

import (
	"errors"
	"regexp"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	first, err := GenerateCode()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateCode()
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !regexp.MustCompile(`^[A-Za-z0-9_-]{22}$`).MatchString(first) {
		t.Fatalf("unexpected code %q", first)
	}
}

func TestGenerateUniqueRetries(t *testing.T) {
	tries := 0
	code, err := GenerateUnique(func(string) (bool, error) { tries++; return tries == 2, nil })
	if err != nil || code == "" || tries != 2 {
		t.Fatalf("code=%q tries=%d err=%v", code, tries, err)
	}
	_, err = GenerateUnique(func(string) (bool, error) { return false, nil })
	if !errors.Is(err, ErrCollisionLimit) {
		t.Fatalf("err=%v", err)
	}
}
