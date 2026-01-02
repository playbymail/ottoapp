--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateSession creates a session for the user.
--
-- name: CreateSession :exec
INSERT INTO sessions (session_id, csrf, user_id, expires_at, created_at, updated_at)
VALUES (:session_id, :csrf, :user_id, :expires_at, :created_at, :updated_at);

-- name: ReadSessionData :one
SELECT sessions.csrf,
       users.user_id,
       users.handle,
       users.is_active,
       users.is_admin,
       users.is_gm,
       users.is_guest,
       users.is_player,
       users.is_service,
       users.is_sysop,
       users.is_user
FROM sessions,
     users
WHERE sessions.session_id = :session_id
  AND sessions.expires_at > :expires_at
  AND users.user_id = sessions.user_id;

-- DeleteSession deletes a session.
--
-- name: DeleteSession :exec
DELETE
FROM sessions
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
DELETE
FROM sessions
WHERE expires_at <= :now_utc;

-- UpdateSession updates a session.
--
-- name: UpdateSession :exec
UPDATE sessions
SET expires_at = :expires_at,
    updated_at = :updated_at
WHERE session_id = :session_id;
