package access

import "context"

// EmailSender is intentionally provider-neutral. Implementations must return
// only a safe delivery category; provider response bodies and credentials never
// cross this boundary.
type EmailSender interface {
	SendStatusAccessEmail(context.Context, StatusAccessEmail) DeliveryResult
}

type StatusAccessEmail struct {
	To              string
	VerificationURL string
	TemplateAlias   string
}

type DeliveryResult string

const (
	DeliveryAccepted         DeliveryResult = "accepted"
	DeliveryTemporaryFailure DeliveryResult = "temporary_failure"
	DeliveryRejected         DeliveryResult = "rejected"
	DeliverySuppressed       DeliveryResult = "suppressed"
)

type FakeEmailSender struct {
	Result DeliveryResult
	Sent   []StatusAccessEmail
}

func (f *FakeEmailSender) SendStatusAccessEmail(_ context.Context, email StatusAccessEmail) DeliveryResult {
	f.Sent = append(f.Sent, email)
	return f.Result
}

// ProviderConfig is deliberately small. SES or another provider can later be
// added without exposing provider-specific details to callers.
type ProviderConfig struct {
	Provider              string
	FromAddress           string
	FromName              string
	MessageStream         string
	PostmarkServerToken   string
	PostmarkTemplateAlias string
}

func NewEmailSender(cfg ProviderConfig) (EmailSender, error) {
	switch cfg.Provider {
	case "fake":
		return &FakeEmailSender{Result: DeliveryAccepted}, nil
	case "postmark":
		return NewPostmarkSender(PostmarkConfig{
			FromAddress:   cfg.FromAddress,
			FromName:      cfg.FromName,
			MessageStream: cfg.MessageStream,
			ServerToken:   cfg.PostmarkServerToken,
			TemplateAlias: cfg.PostmarkTemplateAlias,
		})
	default:
		return nil, ErrUnsupportedProvider
	}
}
