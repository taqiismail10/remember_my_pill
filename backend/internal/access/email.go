package access

import "context"

// EmailSender is intentionally provider-neutral. B3A supplies only this
// contract and a test fake; no external provider or real email is enabled.
type EmailSender interface {
	SendStatusAccessEmail(context.Context, StatusAccessEmail) DeliveryResult
}

type StatusAccessEmail struct {
	To                string
	VerificationToken string
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
