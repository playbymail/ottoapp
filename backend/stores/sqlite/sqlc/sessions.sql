--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateSession creates a session for the user.
--
-- name: CreateSession :exec
INSERT INTO sessions (session_id, csrf, user_id, expires_at)
VALUES(:session_id, :csrf, :user_id, :expires_at);

-- GetSession returns the session tied to an id.
--
-- name: GetSession :one
SELECT csrf, user_id, expires_at
FROM sessions
WHERE session_id = :session_id;

