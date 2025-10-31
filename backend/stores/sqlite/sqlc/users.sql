--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateUser creates a new user and returns its id.
-- The email must be lowercase and unique.
-- Timezone is the user's timezone. Use UTC if unknown.
-- The password is stored as a bcrypt hash.
--
-- name: CreateUser :one
INSERT INTO users (handle,
                   email,
                   timezone)
VALUES (:handle,
        :email,
        :timezone)
RETURNING user_id;

-- GetUserByID returns the user with the given id.
--
-- name: GetUser :one
SELECT handle,
       email,
       timezone,
       created_at,
       updated_at
FROM users
WHERE user_id = :user_id;

-- GetUserIDByEmail returns the id of the user with the given email address.
--
-- name: GetUserIDByEmail :one
SELECT user_id
FROM users
WHERE email = :email;

-- GetUserIDByHandle returns the id of the user with the given handle.
--
-- name: GetUserIDByHandle :one
SELECT user_id
FROM users
WHERE handle = :handle;

-- UpdateUser updates the given user.
--
-- name: UpdateUser :exec
UPDATE users
SET email      = :email,
    handle     = :handle,
    timezone   = :timezone,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;

-- UpdateUserEmail updates the email for the given user.
--
-- name: UpdateUserEmail :exec
UPDATE users
SET email      = :email,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;

-- UpdateUserHandle updates the handle for the given user.
--
-- name: UpdateUserHandle :exec
UPDATE users
SET handle     = :handle,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;

-- UpdateUserTimezone updates the timezone for the given user.
--
-- name: UpdateUserTimezone :exec
UPDATE users
SET timezone   = :timezone,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;

