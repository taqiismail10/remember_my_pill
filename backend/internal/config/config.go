package config

import (
	"fmt"
	"net/netip"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var consentVersionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,29}$`)

type Config struct {
	DatabaseURL                  string
	AllowedOrigin                string
	Port                         string
	ConsentVersion               string
	TrustedProxies               []netip.Prefix
	EmailProvider                string
	EmailFromAddress             string
	EmailFromName                string
	StatusAccessBaseURL          string
	PostmarkServerToken          string
	PostmarkMessageStream        string
	PostmarkStatusAccessTemplate string
	PilotLegalContentApproved    bool
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:                  os.Getenv("DATABASE_URL"),
		AllowedOrigin:                os.Getenv("ALLOWED_ORIGIN"),
		Port:                         os.Getenv("PORT"),
		ConsentVersion:               os.Getenv("CONSENT_VERSION"),
		EmailProvider:                strings.ToLower(strings.TrimSpace(os.Getenv("EMAIL_PROVIDER"))),
		EmailFromAddress:             strings.TrimSpace(os.Getenv("EMAIL_FROM_ADDRESS")),
		EmailFromName:                strings.TrimSpace(os.Getenv("EMAIL_FROM_NAME")),
		StatusAccessBaseURL:          strings.TrimSpace(os.Getenv("STATUS_ACCESS_BASE_URL")),
		PostmarkServerToken:          strings.TrimSpace(os.Getenv("POSTMARK_SERVER_TOKEN")),
		PostmarkMessageStream:        strings.TrimSpace(os.Getenv("POSTMARK_MESSAGE_STREAM")),
		PostmarkStatusAccessTemplate: strings.TrimSpace(os.Getenv("POSTMARK_STATUS_ACCESS_TEMPLATE")),
	}
	legalApproved, err := parseBoolean("PILOT_LEGAL_CONTENT_APPROVED", os.Getenv("PILOT_LEGAL_CONTENT_APPROVED"))
	if err != nil {
		return Config{}, err
	}
	cfg.PilotLegalContentApproved = legalApproved
	if cfg.EmailProvider == "" {
		cfg.EmailProvider = "fake"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if !ValidAllowedOrigin(cfg.AllowedOrigin) {
		return Config{}, fmt.Errorf("ALLOWED_ORIGIN must be an absolute http(s) origin without path")
	}
	if cfg.ConsentVersion != "" && !consentVersionPattern.MatchString(cfg.ConsentVersion) {
		return Config{}, fmt.Errorf("CONSENT_VERSION must be a safe policy identifier")
	}
	proxies, err := parseTrustedProxies(os.Getenv("TRUSTED_PROXY_CIDRS"))
	if err != nil {
		return Config{}, err
	}
	cfg.TrustedProxies = proxies
	if err := validateEmailProvider(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parseBoolean(name, value string) (bool, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "false" {
		return false, nil
	}
	if value == "true" {
		return true, nil
	}
	return false, fmt.Errorf("%s must be true or false", name)
}

func validateEmailProvider(cfg Config) error {
	if cfg.EmailProvider == "fake" {
		return nil
	}
	if cfg.EmailProvider != "postmark" {
		return fmt.Errorf("EMAIL_PROVIDER must be fake or postmark")
	}
	if !validEmailAddress(cfg.EmailFromAddress) || cfg.EmailFromName == "" || strings.ContainsAny(cfg.EmailFromName, "\r\n") || cfg.PostmarkServerToken == "" || !validMessageStream(cfg.PostmarkMessageStream) || !validStatusAccessBaseURL(cfg.StatusAccessBaseURL) {
		return fmt.Errorf("postmark email configuration is incomplete or invalid")
	}
	return nil
}

func validEmailAddress(value string) bool {
	return len(value) <= 255 && strings.Count(value, "@") == 1 && !strings.ContainsAny(value, " \t\r\n") && strings.Contains(strings.Split(value, "@")[1], ".")
}

func validMessageStream(value string) bool {
	if value == "" || len(value) > 80 {
		return false
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-') {
			return false
		}
	}
	return true
}

func validStatusAccessBaseURL(value string) bool {
	u, err := url.ParseRequestURI(value)
	return err == nil && u.Scheme == "https" && u.Host != "" && u.RawQuery == "" && u.Fragment == ""
}

func ValidAllowedOrigin(value string) bool {
	u, err := url.ParseRequestURI(value)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.Path == "" && u.RawQuery == "" && u.Fragment == ""
}

func parseTrustedProxies(value string) ([]netip.Prefix, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	proxies := make([]netip.Prefix, 0, len(parts))
	for _, part := range parts {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("TRUSTED_PROXY_CIDRS contains an invalid CIDR")
		}
		proxies = append(proxies, prefix)
	}
	return proxies, nil
}
