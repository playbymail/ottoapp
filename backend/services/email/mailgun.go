// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package email

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// mailgunService implements Service using Mailgun's HTTP API.
//
// API docs:
//
//	https://documentation.mailgun.com/en/latest/api-sending.html#sending
type mailgunService struct {
	client  *http.Client
	apiBase string // e.g. https://api.mailgun.net/v3
	domain  string // e.g. ottomap.playbymailgames.com
	apiKey  string
	from    string // default From: address
	replyTo string // default Reply-To
}

// NewMailgun returns a Service that sends email via Mailgun's HTTP API.
//
// If apiBase is empty, it defaults to https://api.mailgun.net/v3.
func NewMailgun(apiBase, domain, apiKey, from, replyTo string) (Service, error) {
	if apiBase == "" {
		apiBase = "https://api.mailgun.net/v3"
	}
	if from == "" {
		return nil, fmt.Errorf("from is required")
	}
	return &mailgunService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiBase: apiBase,
		domain:  domain,
		apiKey:  apiKey,
		from:    from,
		replyTo: replyTo,
	}, nil
}

// apiURL builds the full Mailgun messages endpoint URL.
func (s *mailgunService) apiURL() string {
	// Example: https://api.mailgun.net/v3/mg.playbymailgames.com/messages
	return fmt.Sprintf("%s/%s/messages", strings.TrimRight(s.apiBase, "/"), s.domain)
}

// Send sends a generic email message via Mailgun.
func (s *mailgunService) Send(ctx context.Context, msg Message) error {
	form := url.Values{}
	form.Set("from", s.from)
	if msg.ReplyTo != "" { // prefer per-message ReplyTo; fall back to service default.
		form.Set("h:Reply-To", msg.ReplyTo)
	} else if s.replyTo != "" {
		form.Set("h:Reply-To", s.replyTo)
	}
	form.Set("to", msg.To)
	form.Set("subject", msg.Subject)
	if msg.Text != "" {
		form.Set("text", msg.Text)
	}
	if msg.HTML != "" {
		form.Set("html", msg.HTML)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL(), strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Mailgun uses HTTP Basic auth: username "api", password is the API key.
	req.SetBasicAuth("api", s.apiKey)

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	// Consume body for logging in case of error, but avoid large buffering.
	if res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return fmt.Errorf("email provider error: status=%d body=%s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

// SendWelcome sends an account-creation / welcome email.
func (s *mailgunService) SendWelcome(ctx context.Context, to string, name, game string, clan int, secret string) error {
	subject := fmt.Sprintf("%s: %04d: Welcome to OttoMap", game, clan)
	text := fmt.Sprintf(`
Hi %s (Clan %04d)!

Your OttoMap account has been created.
The URL for the web server is https://ottomap.playbymailgames.com/.

You will need to use your email address to log in.

Your password is "%s" (without the quotes).

The server is a work in progress (there are plans, big plans, but I'm not much of a web developer).
You can view your maps (right now that's limited to just the map key) and download them.

I don't want OttoMap to send too many e-mails, so announcements will be made on the Discord server (https://discord.gg/8v2pWUs2Pg).

If there are problems with the server, I may have to e-mail your maps.
I'll post on Discord if that happens.

If you ever want to quit using OttoMap, please send the GM an email telling him that you're opting out of OttoMap.
You'll be removed from the report uploads and your account (with all reports and maps) will be deleted.

– OttoMap

Note: this email was sent from an unmonitored address.
I'm working on fixing that but it's a low priority.
Please contact me on Discord or Gmail if you need anything.`, name, clan, secret)

	html := fmt.Sprintf(`
<p>
Hi %s (Clan %04d)!
</p>
<p>
Your OttoMap account has been created.
The URL for the web server is https://ottomap.playbymailgames.com/.
</p>
<p>
You will need to use your email address to log in.
</p>
Your password is "%s" (without the quotes).
</p>
<p>
The server is a work in progress (there are plans, big plans, but I'm not much of a web developer).
You can view your maps (right now that's limited to just the map key) and download them.
</p>
<p>
I don't want OttoMap to send too many e-mails, so announcements will be made on the Discord server (https://discord.gg/8v2pWUs2Pg).
</p>
<p>
If there are problems with the server, I may have to e-mail your maps.
I'll post on Discord if that happens.
</p>
<p>
If you ever want to quit using OttoMap, please send the GM an email telling him that you're opting out of OttoMap.
You'll be removed from the report uploads and your account (with all reports and maps) will be deleted.
</p>
<p>
– OttoMap
</p>
<p>
Note: this email was sent from an unmonitored address.
I'm working on fixing that but it's a low priority.
Please contact me on Discord or Gmail if you need anything.
</p>`, name, clan, secret)

	return s.Send(ctx, Message{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    html,
	})
}

// SendTurnUploaded notifies a player that turn reports have been uploaded.
func (s *mailgunService) SendTurnUploaded(ctx context.Context, to, name, game string, turn int) error {
	subject := fmt.Sprintf("[%s] Turn %d reports uploaded", game, turn)
	text := fmt.Sprintf(
		"Hi %s,\n\nTurn %d reports for %s have been uploaded to OttoMap. You can now review your report and prepare your orders.\n\n– OttoMap",
		name, turn, game,
	)
	html := fmt.Sprintf(`
        <p>Hi %s,</p>
        <p>Turn <strong>%d</strong> reports for <strong>%s</strong> have been uploaded to OttoMap.</p>
        <p>You can now review your report and prepare your orders.</p>
        <p>– OttoMap</p>
    `, name, turn, game)

	return s.Send(ctx, Message{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    html,
	})
}
