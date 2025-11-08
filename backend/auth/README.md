# Auth service

## Timing Attacks

We pulled the constant timing logic out of the authentication routes because it impacted test performance.

All authentication routes have code similar to the following to cause a constant delay on failed attempts.

```go
// handler
start := time.Now()

id, err := authSvc.AuthenticateWithEmailSecret(email, pass)

// ensure minimum duration
min := 200 * time.Millisecond
if d := time.Since(start); d < min {
    time.Sleep(min - d)
}
```