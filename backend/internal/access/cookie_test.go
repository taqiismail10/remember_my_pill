package access

import (
	"net/http"
	"testing"
	"time"
)

func TestStatusSessionCookieProductionAttributes(t *testing.T) {
	now := time.Date(2026, time.August, 13, 12, 0, 0, 0, time.UTC)
	cookie := StatusSessionCookie("opaque-session", true, now)
	if cookie.Name != ProductionStatusCookieName || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" || cookie.Domain != "" {
		t.Fatalf("unexpected cookie attributes: %+v", *cookie)
	}
	if cookie.MaxAge != int(BrowserSessionTTL.Seconds()) || !cookie.Expires.Equal(now.Add(BrowserSessionTTL)) {
		t.Fatalf("unexpected lifetime: max-age=%d expires=%s", cookie.MaxAge, cookie.Expires)
	}
}

func TestStatusSessionCookieDevelopmentAndClearAttributes(t *testing.T) {
	cookie := StatusSessionCookie("opaque-session", false, time.Now())
	if cookie.Name != DevelopmentStatusCookieName || cookie.Secure || !cookie.HttpOnly || cookie.Domain != "" {
		t.Fatalf("unexpected development cookie attributes: %+v", *cookie)
	}
	cleared := ClearStatusSessionCookie(true)
	if cleared.Name != ProductionStatusCookieName || cleared.MaxAge != -1 || !cleared.Expires.Equal(time.Unix(1, 0)) || !cleared.Secure || !cleared.HttpOnly || cleared.Path != "/" || cleared.Domain != "" {
		t.Fatalf("unexpected clear cookie attributes: %+v", *cleared)
	}
}
