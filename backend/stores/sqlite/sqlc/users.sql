--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- UpsertUser creates a new user and returns its id.
-- The email must be lowercase and unique.
-- Timezone is the user's timezone. Use UTC if unknown.
-- The password is stored as a bcrypt hash.
--
-- name: UpsertUser :one
INSERT INTO users (handle,
                   email,
                   username,
                   timezone,
                   created_at,
                   updated_at)
VALUES (:handle,
        :email,
        :username,
        :timezone,
        :created_at,
        :updated_at)
ON CONFLICT (handle) DO UPDATE
    SET email      = excluded.email,
        username   = excluded.username,
        timezone   = excluded.timezone,
        updated_at = excluded.updated_at
RETURNING user_id;

-- GetUserByID returns the user with the given id.
--
-- name: GetUserByID :one
SELECT user_id,
       username,
       email,
       handle,
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
       handle,
       timezone,
       created_at,
       updated_at
FROM users
WHERE email = :email;

-- GetUserByHandle returns the user with the given handle.
--
-- name: GetUserByHandle :one
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
WHERE handle = :handle;

-- GetUserByUsername returns the user with the given username.
--
-- name: GetUserByUsername :one
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
WHERE username = :username;

-- name: GetUserIDByEmail :one
SELECT user_id
FROM users
WHERE email = :email;

-- name: GetUserIDByHandle :one
SELECT user_id
FROM users
WHERE handle = :handle;

-- name: GetUserIDByUsername :one
SELECT user_id
FROM users
WHERE username = :username;

-- GetAllUsers returns all users.
--
-- name: GetAllUsers :many
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
ORDER BY username;

-- name: ListUsersVisibleToActor :many
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
WHERE :actor_id = 1
  AND :page_size = 1
  AND :page_num = 1
ORDER BY username;