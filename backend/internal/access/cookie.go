package access

import (
	"net/http"
	"time"
)

const (
	ProductionStatusCookieName  = "__Host-rmp-status"
	DevelopmentStatusCookieName = "rmp-status-dev"
)

func StatusSessionCookie(rawSessionID string, secure bool, now time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     statusSessionCookieName(secure),
		Value:    rawSessionID,
		Path:     "/",
		MaxAge:   int(BrowserSessionTTL.Seconds()),
		Expires:  now.Add(BrowserSessionTTL),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func ClearStatusSessionCookie(secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     statusSessionCookieName(secure),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func statusSessionCookieName(secure bool) string {
	if secure {
		return ProductionStatusCookieName
	}
	// __Host- cookies require Secure. The development name makes the deliberate
	// local HTTP exception visible and cannot be mistaken for production scope.
	return DevelopmentStatusCookieName
}
