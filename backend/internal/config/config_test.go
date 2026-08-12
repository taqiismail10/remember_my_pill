package config

import "testing"

func TestValidAllowedOrigin(t *testing.T) {
	for _, value := range []string{"", "localhost:3000", "ftp://example.test", "http://example.test/path", "https://example.test"} {
		want := value == "https://example.test"
		if got := ValidAllowedOrigin(value); got != want {
			t.Fatalf("%q = %v", value, got)
		}
	}
}

func TestParseTrustedProxies(t *testing.T) {
	proxies, err := parseTrustedProxies("127.0.0.1/32, 10.0.0.0/8")
	if err != nil || len(proxies) != 2 {
		t.Fatalf("proxies=%v err=%v", proxies, err)
	}
	if _, err := parseTrustedProxies("not-a-cidr"); err == nil {
		t.Fatal("expected invalid CIDR error")
	}
}

func TestValidateEmailProvider(t *testing.T) {
	postmark := Config{EmailProvider: "postmark", EmailFromAddress: "status@example.test", EmailFromName: "Remember My Pill", StatusAccessBaseURL: "https://remembermypill.example", PostmarkServerToken: "test-token", PostmarkMessageStream: "outbound"}
	if err := validateEmailProvider(postmark); err != nil {
		t.Fatalf("valid postmark config: %v", err)
	}
	for _, cfg := range []Config{{EmailProvider: "other"}, {EmailProvider: "postmark"}, {EmailProvider: "postmark", EmailFromAddress: "status@example.test", EmailFromName: "RMP", StatusAccessBaseURL: "http://rmp.example", PostmarkServerToken: "test", PostmarkMessageStream: "outbound"}} {
		if err := validateEmailProvider(cfg); err == nil {
			t.Fatalf("expected invalid provider config: %+v", cfg)
		}
	}
}

func TestLoadDefaultsToFakeWithoutPostmarkCredentials(t *testing.T) {
	for _, key := range []string{"DATABASE_URL", "ALLOWED_ORIGIN", "PORT", "CONSENT_VERSION", "TRUSTED_PROXY_CIDRS", "EMAIL_PROVIDER", "EMAIL_FROM_ADDRESS", "EMAIL_FROM_NAME", "STATUS_ACCESS_BASE_URL", "POSTMARK_SERVER_TOKEN", "POSTMARK_MESSAGE_STREAM", "POSTMARK_STATUS_ACCESS_TEMPLATE", "PILOT_LEGAL_CONTENT_APPROVED", "RMP_ENVIRONMENT", "STATUS_SESSION_COOKIE_SECURE"} {
		t.Setenv(key, "")
	}
	t.Setenv("DATABASE_URL", "postgres://example.test/db")
	t.Setenv("ALLOWED_ORIGIN", "http://example.test")
	if cfg, err := Load(); err != nil || cfg.EmailProvider != "fake" || !cfg.StatusSessionCookieSecure || cfg.RuntimeEnvironment != "production" {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
}

func TestStatusSessionCookieSecurityConfiguration(t *testing.T) {
	for _, key := range []string{"DATABASE_URL", "ALLOWED_ORIGIN", "RMP_ENVIRONMENT", "STATUS_SESSION_COOKIE_SECURE"} {
		t.Setenv(key, "")
	}
	t.Setenv("DATABASE_URL", "postgres://example.test/db")
	t.Setenv("ALLOWED_ORIGIN", "http://example.test")
	t.Setenv("RMP_ENVIRONMENT", "development")
	t.Setenv("STATUS_SESSION_COOKIE_SECURE", "false")
	if cfg, err := Load(); err != nil || cfg.StatusSessionCookieSecure {
		t.Fatalf("development cfg=%+v err=%v", cfg, err)
	}

	t.Setenv("RMP_ENVIRONMENT", "production")
	if _, err := Load(); err == nil {
		t.Fatal("expected insecure production cookie configuration to fail")
	}
}

func TestParseLegalContentApproval(t *testing.T) {
	if value, err := parseBoolean("PILOT_LEGAL_CONTENT_APPROVED", "true"); err != nil || !value {
		t.Fatalf("true=%v err=%v", value, err)
	}
	if value, err := parseBoolean("PILOT_LEGAL_CONTENT_APPROVED", "false"); err != nil || value {
		t.Fatalf("false=%v err=%v", value, err)
	}
	if _, err := parseBoolean("PILOT_LEGAL_CONTENT_APPROVED", "yes"); err == nil {
		t.Fatal("expected invalid boolean")
	}
}
