-- name: GetSessionData :many
SELECT s.csrf,
       u.user_id,
       u.handle,
       ur.role_id
FROM sessions s,
     users u,
     user_roles ur
WHERE s.session_id = :session_id
  AND s.expires_at > :expires_at
  AND u.user_id = s.user_id
  AND ur.user_id = u.user_id
;
