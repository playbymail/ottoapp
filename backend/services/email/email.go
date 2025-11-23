package email

import "context"

// Service defines the operations for sending transactional emails from OttoMap.
//
// Implementations may use any provider (Mailgun, Postmark, SES, SMTP, etc.).
// Production code should depend on this interface, not on a specific provider.
type Service interface {
	// Send sends a generic email message.
	Send(ctx context.Context, msg Message) error

	// SendWelcome sends an account-creation / welcome email.
	SendWelcome(ctx context.Context, to string, name, game string, clan int, secret string) error

	// SendTurnUploaded notifies a player that turn reports were uploaded.
	SendTurnUploaded(ctx context.Context, to, name, game string, turn int) error
}

// Message represents a generic outbound email.
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
	ReplyTo string // optional Reply-To header
}
