--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateSession creates a session for the user.
--
-- name: CreateSession :exec
INSERT INTO sessions (session_id, csrf, user_id, expires_at)
VALUES (:session_id, :csrf, :user_id, :expires_at);

-- DeleteSession deletes a session.
--
-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_id = :session_id;

-- GetSession returns the session tied to an id.
--
-- name: GetSession :one
SELECT csrf, user_id, expires_at
FROM sessions
WHERE session_id = :session_id;

-- ReapSessions deletes expired sessions.
--
-- name: ReapSessions :exec
DELETE FROM sessions
WHERE expires_at <= :now_utc;

-- UpdateSession updates a session.
--
-- name: UpdateSession :exec
UPDATE sessions
SET expires_at = :expires_at
WHERE session_id = :session_id;
