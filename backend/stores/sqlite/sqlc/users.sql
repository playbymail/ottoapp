-- name: CreateUser :one
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
RETURNING user_id;

-- name: ReadEmailByUserId :one
SELECT email
FROM users
WHERE user_id = :user_id;

-- name: ReadHandleByUserId :one
SELECT handle
FROM users
WHERE user_id = :user_id;

-- name: ReadUserByUserId :one
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
WHERE user_id = :user_id;

-- name: ReadUserIdByEmail :one
SELECT user_id
FROM users
WHERE email = :email;

-- ReadUserIdByHandle returns the id of the user with the given handle.
--
-- name: ReadUserIdByHandle :one
SELECT user_id
FROM users
WHERE handle = :handle;

-- ReadUserIdByUsername returns the id of the user with the given username.
--
-- name: ReadUserIdByUsername :one
SELECT user_id
FROM users
WHERE username = :username;

-- name: ReadUsers :many
SELECT user_id,
       username,
       email,
       handle,
       timezone,
       created_at,
       updated_at
FROM users
ORDER BY username;

-- name: ReadUsersVisibleToActor :many
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

-- name: UpdateEmailByUserId :exec
UPDATE users
SET email      = LOWER(:email),
    updated_at = :updated_at
WHERE user_id = :user_id;

-- name: UpdateHandleByUserId :exec
UPDATE users
SET handle     = LOWER(:handle),
    updated_at = :updated_at
WHERE user_id = :user_id;

-- name: UpdateTimezoneByUserId :exec
UPDATE users
SET timezone   = :timezone,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- name: UpdateUserByUserId :exec
UPDATE users
SET email      = LOWER(:email),
    handle     = LOWER(:handle),
    timezone   = :timezone,
    username   = :username,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- name: UpdateUsernameByUserId :exec
UPDATE users
SET username   = :username,
    updated_at = :updated_at
WHERE user_id = :user_id;
