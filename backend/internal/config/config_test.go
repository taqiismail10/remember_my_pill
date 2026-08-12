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
