# OttoMap Email Service — Mailgun Version

## Layout

```text
backend/
  services/
    email/
      email.go      # Service interface + Message type
      mailgun.go    # Mailgun-backed implementation + domain helpers
```

## Configuration

The Mailgun implementation expects these environment variables (or equivalent
config values) to be provided at startup:

- `MAILGUN_DOMAIN`   — e.g. `mg.playbymailgames.com`
- `MAILGUN_API_KEY`  — your Mailgun private API key
- `EMAIL_FROM`       — e.g. `OttoMap <ottomap@mg.playbymailgames.com>`
- `MAILGUN_API_BASE` — optional; defaults to `https://api.mailgun.net/v3`

## Usage Sketch

```go
import "github.com/your/module/backend/services/email"

func main() {
    domain := os.Getenv("MAILGUN_DOMAIN")
    apiKey := os.Getenv("MAILGUN_API_KEY")
    from   := os.Getenv("EMAIL_FROM")
    apiBase := os.Getenv("MAILGUN_API_BASE") // optional

    emailSvc := email.NewMailgun(apiBase, domain, apiKey, from)

    // On registration:
    if err := emailSvc.SendWelcome(ctx, user.Email, user.Name); err != nil {
        log.Printf("welcome email failed: %v", err)
    }

    // After turn processing:
    if err := emailSvc.SendTurnUploaded(ctx, player.Email, player.Name, gameID, turnNo); err != nil {
        log.Printf("turn email failed: %v", err)
    }
}
```

You can extend the `Service` interface with more helpers as new email types are added.
