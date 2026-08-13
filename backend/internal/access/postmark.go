package access

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const postmarkSendWithTemplateURL = "https://api.postmarkapp.com/email/withTemplate"

var (
	ErrUnsupportedProvider   = errors.New("unsupported email provider")
	ErrPostmarkConfiguration = errors.New("invalid postmark configuration")
)

// PostmarkConfig is adapter-only and must not be used by access-domain
// callers. The token is supplied from environment configuration only.
type PostmarkConfig struct {
	FromAddress   string
	FromName      string
	MessageStream string
	ServerToken   string
	TemplateAlias string
	Endpoint      string // Test seam; production uses Postmark's documented endpoint.
	HTTPClient    *http.Client
}

type PostmarkSender struct {
	fromAddress   string
	fromName      string
	messageStream string
	serverToken   string
	templateAlias string
	endpoint      string
	client        *http.Client
}

// NewPostmarkSender validates configuration without sending email. A template
// alias is intentionally optional here because final production copy is not
// approved; SendStatusAccessEmail rejects without one and performs no request.
func NewPostmarkSender(cfg PostmarkConfig) (*PostmarkSender, error) {
	if !validEmail(cfg.FromAddress) || strings.TrimSpace(cfg.FromName) == "" ||
		strings.ContainsAny(cfg.FromName, "\r\n") || strings.TrimSpace(cfg.ServerToken) == "" ||
		!validMessageStream(cfg.MessageStream) {
		return nil, ErrPostmarkConfiguration
	}
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = postmarkSendWithTemplateURL
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, ErrPostmarkConfiguration
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &PostmarkSender{
		fromAddress: cfg.FromAddress, fromName: cfg.FromName, messageStream: cfg.MessageStream,
		serverToken: cfg.ServerToken, templateAlias: cfg.TemplateAlias, endpoint: endpoint, client: client,
	}, nil
}

// SendStatusAccessEmail sends only approved Postmark templates. It never
// logs, stores, or returns recipient, URL, token, or provider error content.
func (s *PostmarkSender) SendStatusAccessEmail(ctx context.Context, email StatusAccessEmail) DeliveryResult {
	template := s.templateAlias
	if template == "" || (email.TemplateAlias != "" && email.TemplateAlias != template) || !validEmail(email.To) || !validVerificationURL(email.VerificationURL) {
		return DeliveryRejected
	}
	payload := struct {
		From          string         `json:"From"`
		To            string         `json:"To"`
		TemplateAlias string         `json:"TemplateAlias"`
		TemplateModel map[string]any `json:"TemplateModel"`
		MessageStream string         `json:"MessageStream"`
		Tag           string         `json:"Tag"`
	}{
		From:          fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress),
		To:            email.To,
		TemplateAlias: template,
		TemplateModel: map[string]any{"verification_url": email.VerificationURL},
		MessageStream: s.messageStream,
		Tag:           "status-access",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return DeliveryRejected
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return DeliveryRejected
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.serverToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return DeliveryTemporaryFailure
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return DeliveryTemporaryFailure
	}
	limited := io.LimitReader(resp.Body, 64*1024)
	var response struct {
		ErrorCode int    `json:"ErrorCode"`
		MessageID string `json:"MessageID"`
	}
	if err := json.NewDecoder(limited).Decode(&response); err != nil {
		return DeliveryTemporaryFailure
	}
	if response.ErrorCode == 0 && resp.StatusCode >= 200 && resp.StatusCode < 300 && response.MessageID != "" {
		return DeliveryAccepted
	}
	// Postmark uses error 406 for an inactive/suppressed recipient.
	if response.ErrorCode == 406 {
		return DeliverySuppressed
	}
	return DeliveryRejected
}

// BuildVerificationURL puts a short-lived verification value only in a URL
// fragment. Fragments are not sent to the server in the initial navigation;
// the browser verification route immediately POSTs it and clears history.
func BuildVerificationURL(baseURL, verificationToken string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil || base.Scheme != "https" || base.Host == "" || base.RawQuery != "" || base.Fragment != "" || verificationToken == "" {
		return "", errors.New("invalid status access URL configuration")
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/waitlist/verify"
	base.RawQuery = ""
	base.Fragment = url.Values{"v": []string{verificationToken}}.Encode()
	return base.String(), nil
}

func validVerificationURL(value string) bool {
	u, err := url.Parse(value)
	fragment, fragmentErr := url.ParseQuery(u.Fragment)
	return err == nil && fragmentErr == nil && u.Scheme == "https" && u.Host != "" && u.Path == "/waitlist/verify" && u.RawQuery == "" && fragment.Get("v") != ""
}

func validEmail(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) <= 255 && strings.Count(value, "@") == 1 && !strings.ContainsAny(value, " \t\r\n") && strings.Contains(strings.Split(value, "@")[1], ".")
}

func validMessageStream(value string) bool {
	value = strings.TrimSpace(value)
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
