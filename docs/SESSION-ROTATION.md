Here’s fairly concrete pseudo-code (Go-flavored) for **cookie-based opaque sessions** with:

* sliding idle expiry (extend on activity)
* absolute expiry (hard cap)
* “touch interval” (avoid DB write every request)
* token rotation with grace window (avoids parallel-request breakage)
* revocation support

### Data model (conceptual)

```text
Session {
  id_hash            // hash(token) primary key (store hash, not raw)
  user_id
  created_at
  absolute_expires_at
  last_seen_at
  rotated_at
  prev_id_hash       // optional (for grace window)
  prev_expires_at    // optional
  revoked_at         // nullable
  user_agent_sig     // optional soft bind (hash of UA)
}
```

### Config

```text
IDLE_TTL            = 45 days         // sliding
ABSOLUTE_TTL        = 180 days        // hard cap from created_at
TOUCH_INTERVAL      = 30 minutes      // only update last_seen at most this often
ROTATE_INTERVAL     = 14 days
ROTATION_GRACE      = 5 minutes       // accept previous token briefly
COOKIE_NAME         = "sess"
```

---

## Middleware pseudo-code

```go
func SessionMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // 1) Extract token from cookie
    token := readCookie(r, COOKIE_NAME)
    if token == "" {
      // anonymous request
      next.ServeHTTP(w, r)
      return
    }

    // 2) Hash token for lookup (store only hashes in DB)
    tokenHash := hashToken(token)

    // 3) Lookup session:
    //    - match on id_hash OR (prev_id_hash within grace)
    now := timeNow()
    sess, matchKind := dbFindSessionByHash(tokenHash, now)
    // matchKind: "current" | "previous" | "none"

    if sess == nil {
      clearCookie(w, COOKIE_NAME)
      next.ServeHTTP(w, r)
      return
    }

    // 4) Reject revoked or expired sessions
    if sess.revoked_at != nil {
      clearCookie(w, COOKIE_NAME)
      next.ServeHTTP(w, r)
      return
    }

    if now.After(sess.absolute_expires_at) {
      // absolute cap reached
      dbRevokeSession(sess) // or delete
      clearCookie(w, COOKIE_NAME)
      next.ServeHTTP(w, r)
      return
    }

    // Sliding idle expiry (based on last_seen)
    if now.Sub(sess.last_seen_at) > IDLE_TTL {
      dbRevokeSession(sess)
      clearCookie(w, COOKIE_NAME)
      next.ServeHTTP(w, r)
      return
    }

    // 5) Optional: soft-bind checks (don’t overdo it)
    //    If UA differs wildly, you can require re-auth or step-up.
    if sess.user_agent_sig != "" {
      if hashUA(r.UserAgent()) != sess.user_agent_sig {
        // Conservative option: invalidate
        // Safer option: mark for step-up auth (not shown)
        dbRevokeSession(sess)
        clearCookie(w, COOKIE_NAME)
        next.ServeHTTP(w, r)
        return
      }
    }

    // 6) Attach user identity to context for downstream handlers
    // (You can lazy-load user record later; store user_id now)
    ctx := contextWithUserID(r.Context(), sess.user_id)
    r = r.WithContext(ctx)

    // 7) Decide whether to TOUCH last_seen (avoid DB write every request)
    shouldTouch := now.Sub(sess.last_seen_at) >= TOUCH_INTERVAL

    // 8) Decide whether to ROTATE token
    shouldRotate := now.Sub(sess.rotated_at) >= ROTATE_INTERVAL

    if shouldRotate {
      // 8a) Mint new token and hash
      newToken := randomToken()
      newHash  := hashToken(newToken)

      // 8b) Update session record with rotation grace
      // Accept old token briefly to avoid breaking parallel requests
      sess.prev_id_hash     = sess.id_hash
      sess.prev_expires_at  = now.Add(ROTATION_GRACE)

      sess.id_hash          = newHash
      sess.rotated_at       = now
      // Also touch last_seen now (rotation counts as activity)
      sess.last_seen_at     = now

      dbUpdateSessionRotation(sess)

      // 8c) Set new cookie
      setCookie(w, COOKIE_NAME, newToken, /*expires*/ now.Add(IDLE_TTL))
      // Note: cookie expiry is not authoritative; server checks are.

    } else if shouldTouch {
      // 7a) Touch last_seen and maybe extend nothing else
      // You do NOT extend absolute_expires_at.
      dbTouchSessionLastSeen(sess.id_hash, now)
      // Optionally also extend cookie expiry to match sliding window
      // (server remains authoritative)
      setCookie(w, COOKIE_NAME, token, /*expires*/ now.Add(IDLE_TTL))
    }

    // 9) Continue
    next.ServeHTTP(w, r)
  })
}
```

---

## `dbFindSessionByHash` pseudo-code

This is where the “rotation grace window” lives.

```go
func dbFindSessionByHash(tokenHash string, now time.Time) (*Session, string) {
  // 1) Try current token
  sess := query(`
    SELECT * FROM sessions
    WHERE id_hash = ? LIMIT 1
  `, tokenHash)
  if sess != nil {
    return sess, "current"
  }

  // 2) Try previous token within grace window
  sess = query(`
    SELECT * FROM sessions
    WHERE prev_id_hash = ?
      AND prev_expires_at IS NOT NULL
      AND prev_expires_at >= ?
    LIMIT 1
  `, tokenHash, now)
  if sess != nil {
    return sess, "previous"
  }

  return nil, "none"
}
```

---

## Password change / “logout all devices”

This is how you keep long-ish sessions from being scary.

```go
func OnPasswordChanged(userID int) {
  // revoke all sessions for user
  exec(`UPDATE sessions SET revoked_at = NOW() WHERE user_id = ? AND revoked_at IS NULL`, userID)
}
```

---

## Notes you’ll care about (given remote MariaDB latency)

* Reads happen per request; writes happen only on:

  * touch interval boundary (e.g., every 30 minutes per active user)
  * rotation interval (e.g., every 14 days per active user)
* That’s deliberately “DB-write-thrifty.”
* If you want even fewer DB reads, you can use a short-lived signed “session hint” cookie, but that’s complexity you probably don’t need.
