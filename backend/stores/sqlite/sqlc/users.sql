--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateUser creates a new user and returns its id.
-- The email must be lowercase and unique.
-- Timezone is the user's timezone. Use UTC if unknown.
-- The password is stored as a bcrypt hash.
--
-- name: CreateUser :one
INSERT INTO users (username,
                   email,
                   timezone,
                   created_at,
                   updated_at)
VALUES (:username,
        :email,
        :timezone,
        :created_at,
        :updated_at)
RETURNING user_id;

-- GetUserByID returns the user with the given id.
--
-- name: GetUserByID :one
SELECT user_id,
       username,
       email,
       timezone,
       created_at,
       updated_at
FROM users
WHERE user_id = :user_id;

-- GetUserByEmail returns the user with the given email address.
--
-- name: GetUserByEmail :one
SELECT user_id,
       username,
       email,
       timezone,
       created_at,
       updated_at
FROM users
WHERE email = :email;

-- GetUserByUsername returns the user with the given username.
--
-- name: GetUserByUsername :one
SELECT user_id,
       username,
       email,
       timezone,
       created_at,
       updated_at
FROM users
WHERE username = :username;

-- name: GetUserIDByEmail :one
SELECT user_id
FROM users
WHERE email = :email;

-- name: GetUserIDByUsername :one
SELECT user_id
FROM users
WHERE username = :username;

-- UpdateUser updates the given user.
--
-- name: UpdateUser :exec
UPDATE users
SET email      = :email,
    username   = :username,
    timezone   = :timezone,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- UpdateUserEmail updates the email for the given user.
--
-- name: UpdateUserEmail :exec
UPDATE users
SET email      = :email,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- UpdateUserName updates the name for the given user.
--
-- name: UpdateUserName :exec
UPDATE users
SET username   = :username,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- UpdateUserTimezone updates the timezone for the given user.
--
-- name: UpdateUserTimezone :exec
UPDATE users
SET timezone   = :timezone,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- GetAllUsers returns all users.
--
-- name: GetAllUsers :many
SELECT user_id,
       username,
       email,
       timezone,
       created_at,
       updated_at
FROM users
ORDER BY username;

-- name: ListUsersVisibleToActor :many
SELECT user_id,
       username,
       email,
       timezone,
       created_at,
       updated_at
FROM users
WHERE :actor_id = 1
  AND :page_size = 1
  AND :page_num = 1
ORDER BY username;