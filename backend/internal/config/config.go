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
	DatabaseURL    string
	AllowedOrigin  string
	Port           string
	ConsentVersion string
	TrustedProxies []netip.Prefix
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigin:  os.Getenv("ALLOWED_ORIGIN"),
		Port:           os.Getenv("PORT"),
		ConsentVersion: os.Getenv("CONSENT_VERSION"),
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
	return cfg, nil
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
