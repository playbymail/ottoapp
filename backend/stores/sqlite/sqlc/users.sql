-- name: CreateUser :one
INSERT INTO users (handle,
                   username,
                   email,
                   email_opt_in,
                   timezone,
                   is_active,
                   is_admin,
                   is_gm,
                   is_guest,
                   is_player,
                   is_service,
                   is_user,
                   created_at,
                   updated_at)
VALUES (:handle,
        :username,
        :email,
        :email_opt_in,
        :timezone,
        :is_active,
        :is_admin,
        :is_gm,
        :is_guest,
        :is_player,
        :is_service,
        :is_user,
        :created_at,
        :updated_at)
ON CONFLICT (user_id) DO UPDATE
    SET handle       = excluded.handle,
        username     = excluded.username,
        email        = excluded.email,
        email_opt_in = excluded.email_opt_in,
        timezone     = excluded.timezone,
        is_active    = excluded.is_active,
        is_admin     = excluded.is_admin,
        is_gm        = excluded.is_gm,
        is_guest     = excluded.is_guest,
        is_player    = excluded.is_player,
        is_service   = excluded.is_service,
        is_user      = excluded.is_user,
        updated_at   = excluded.updated_at
RETURNING user_id;

-- name: ReadUserByEmail :one
SELECT user_id,
       handle,
       username,
       email,
       email_opt_in,
       timezone,
       is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user,
       created_at,
       updated_at
FROM users
WHERE email = :email;

-- name: ReadUserByHandle :one
SELECT user_id,
       handle,
       username,
       email,
       email_opt_in,
       timezone,
       is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user,
       created_at,
       updated_at
FROM users
WHERE handle = :handle;

-- name: ReadUserByUserId :one
SELECT user_id,
       handle,
       username,
       email,
       email_opt_in,
       timezone,
       is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user,
       created_at,
       updated_at
FROM users
WHERE user_id = :user_id;

-- name: ReadUserRoles :one
SELECT is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user
FROM users
WHERE user_id = :user_id;

-- ReadUserSecret returns the password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: ReadUserSecret :one
SELECT hashed_password
FROM users
WHERE user_id = :user_id;

-- name: ReadUsers :many
SELECT user_id,
       handle,
       username,
       email,
       email_opt_in,
       timezone,
       is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user,
       created_at,
       updated_at
FROM users
ORDER BY username;

-- name: ReadUsersVisibleToActor :many
SELECT user_id,
       handle,
       username,
       email,
       email_opt_in,
       timezone,
       is_active,
       is_admin,
       is_gm,
       is_guest,
       is_player,
       is_service,
       is_sysop,
       is_user,
       created_at,
       updated_at
FROM users
WHERE :actor_id = 1
  AND :page_size = 1
  AND :page_num = 1
ORDER BY username;

-- name: UpdateUser :exec
UPDATE users
SET handle       = :handle,
    username     = :username,
    email        = :email,
    email_opt_in = :email_opt_in,
    timezone     = :timezone,
    is_active    = :is_active,
    is_admin     = :is_admin,
    is_gm        = :is_gm,
    is_guest     = :is_guest,
    is_player    = :is_player,
    is_service   = :is_service,
    is_sysop     = :is_sysop,
    is_user      = :is_user,
    updated_at   = :updated_at
WHERE :user_id = :user_id;

-- name: UpdateEmailByUserId :exec
UPDATE users
SET email        = LOWER(:email),
    email_opt_in = :email_opt_in,
    updated_at   = :updated_at
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

-- name: UpdateUsernameByUserId :exec
UPDATE users
SET username   = :username,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- UpdateUserSecret stores the secret for the user.
-- The password is stored as a bcrypt hash.
--
-- name: UpdateUserSecret :exec
UPDATE users
SET hashed_password    = :hashed_password,
    plaintext_password = :plaintext_password,
    updated_at         = :updated_at
WHERE user_id = :user_id;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = :last_login,
    updated_at = :updated_at
WHERE user_id = :user_id;

-- name: UpdateUserRoles :exec
UPDATE users
SET is_active  = :is_active,
    is_admin   = :is_admin,
    is_gm      = :is_gm,
    is_guest   = :is_guest,
    is_player  = :is_player,
    is_user    = :is_user,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_service = 0
  AND is_sysop = 0;

-- name: UpdateUserActiveRole :exec
UPDATE users
SET is_active  = :has_role,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;

-- name: UpdateUserAdminRole :exec
UPDATE users
SET is_admin   = :has_role,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;

-- name: UpdateUserGMRole :exec
UPDATE users
SET is_gm      = :has_role,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;

-- name: UpdateUserGuestRole :exec
UPDATE users
SET is_active  = 1,
    is_admin   = 0,
    is_gm      = 0,
    is_guest   = :has_role,
    is_player  = 0,
    is_user    = 0,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;

-- name: UpdateUserPlayerRole :exec
UPDATE users
SET is_player  = :has_role,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;

-- name: UpdateUserUserRole :exec
UPDATE users
SET is_user    = :has_role,
    updated_at = :updated_at
WHERE user_id = :user_id
  AND is_sysop = 0;
