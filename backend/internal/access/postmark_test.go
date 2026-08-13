package access

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPostmarkSenderOutcomes(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   DeliveryResult
	}{
		{"accepted", http.StatusOK, `{"ErrorCode":0,"MessageID":"safe-id"}`, DeliveryAccepted},
		{"suppressed", http.StatusUnprocessableEntity, `{"ErrorCode":406}`, DeliverySuppressed},
		{"rejected", http.StatusUnprocessableEntity, `{"ErrorCode":300}`, DeliveryRejected},
		{"temporary server error", http.StatusServiceUnavailable, `{"ErrorCode":100}`, DeliveryTemporaryFailure},
		{"rate limited", http.StatusTooManyRequests, `{"ErrorCode":0}`, DeliveryTemporaryFailure},
		{"malformed", http.StatusOK, `{`, DeliveryTemporaryFailure},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Header.Get("X-Postmark-Server-Token") != "test-server-token" {
					t.Error("missing Postmark authentication header")
				}
				if req.URL.Path != "/email/withTemplate" {
					t.Errorf("path = %q", req.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			sender := testPostmarkSender(t, server.URL+"/email/withTemplate", &http.Client{Timeout: time.Second})
			if got := sender.SendStatusAccessEmail(context.Background(), testEmail()); got != tt.want {
				t.Fatalf("outcome = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPostmarkSenderRejectsUnconfiguredTemplateWithoutNetwork(t *testing.T) {
	sender, err := NewPostmarkSender(PostmarkConfig{FromAddress: "status@example.test", FromName: "RMP", MessageStream: "outbound", ServerToken: "test-server-token", Endpoint: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sender.SendStatusAccessEmail(context.Background(), testEmail()); got != DeliveryRejected {
		t.Fatalf("outcome = %q", got)
	}
}

func TestPostmarkSenderRejectsUnexpectedTemplateAlias(t *testing.T) {
	sender := testPostmarkSender(t, "http://127.0.0.1:1", &http.Client{Timeout: time.Second})
	email := testEmail()
	email.TemplateAlias = "unapproved-template"
	if got := sender.SendStatusAccessEmail(context.Background(), email); got != DeliveryRejected {
		t.Fatalf("outcome = %q", got)
	}
}

func TestPostmarkSenderTimeoutAndNetworkFailureAreTemporary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { time.Sleep(100 * time.Millisecond) }))
	defer server.Close()
	sender := testPostmarkSender(t, server.URL, &http.Client{Timeout: time.Millisecond})
	if got := sender.SendStatusAccessEmail(context.Background(), testEmail()); got != DeliveryTemporaryFailure {
		t.Fatalf("timeout outcome = %q", got)
	}
	sender = testPostmarkSender(t, "http://127.0.0.1:1", &http.Client{Timeout: time.Second})
	if got := sender.SendStatusAccessEmail(context.Background(), testEmail()); got != DeliveryTemporaryFailure {
		t.Fatalf("network outcome = %q", got)
	}
}

func TestPostmarkConfigurationAndProviderSelection(t *testing.T) {
	if _, err := NewPostmarkSender(PostmarkConfig{}); err == nil {
		t.Fatal("expected configuration error")
	}
	if sender, err := NewEmailSender(ProviderConfig{Provider: "fake"}); err != nil || sender == nil {
		t.Fatalf("fake selection: %v", err)
	}
	if sender, err := NewEmailSender(ProviderConfig{Provider: "postmark", FromAddress: "status@example.test", FromName: "RMP", MessageStream: "outbound", PostmarkServerToken: "test-server-token"}); err != nil || sender == nil {
		t.Fatalf("postmark selection: %v", err)
	}
	if _, err := NewEmailSender(ProviderConfig{Provider: "other"}); err == nil {
		t.Fatal("expected provider error")
	}
}

func TestBuildVerificationURL(t *testing.T) {
	value, err := BuildVerificationURL("https://remembermypill.com", "short-lived-token")
	if err != nil || !strings.Contains(value, "/waitlist/verify#v=short-lived-token") || strings.Contains(value, "?") {
		t.Fatalf("url=%q err=%v", value, err)
	}
	if _, err := BuildVerificationURL("http://remembermypill.com", "token"); err == nil {
		t.Fatal("expected HTTPS validation")
	}
}

func testPostmarkSender(t *testing.T, endpoint string, client *http.Client) *PostmarkSender {
	t.Helper()
	sender, err := NewPostmarkSender(PostmarkConfig{FromAddress: "status@example.test", FromName: "RMP", MessageStream: "outbound", ServerToken: "test-server-token", TemplateAlias: "status-access-test", Endpoint: endpoint, HTTPClient: client})
	if err != nil {
		t.Fatal(err)
	}
	return sender
}

func testEmail() StatusAccessEmail {
	return StatusAccessEmail{To: "person@example.test", VerificationURL: "https://remembermypill.example/waitlist/verify#v=short-lived-token"}
}
