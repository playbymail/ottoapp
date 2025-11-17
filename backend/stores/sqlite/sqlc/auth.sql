--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- UpsertUserSecret creates a secrets record for the user.
-- The password is stored as a bcrypt hash.
--
-- name: UpsertUserSecret :exec
INSERT INTO user_secrets (user_id,
                          hashed_password,
                          plaintext_password,
                          last_login,
                          created_at,
                          updated_at)
VALUES (:user_id,
        :hashed_password,
        :plaintext_password,
        :last_login,
        :created_at,
        :updated_at)
ON CONFLICT (user_id) DO UPDATE SET hashed_password    = excluded.hashed_password,
                                    plaintext_password = excluded.plaintext_password,
                                    last_login         = excluded.last_login,
                                    updated_at         = excluded.updated_at;

-- GetUserSecret returns the password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: GetUserSecret :one
SELECT hashed_password
FROM user_secrets
WHERE user_id = :user_id;

-- name: UpdateUserLastLogin :exec
UPDATE user_secrets
SET last_login = :last_login
WHERE user_id = :user_id;

-- UpsertUserRole assigns a role to a user.
--
-- name: UpsertUserRole :exec
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (:user_id, :role_id, :created_at, :updated_at)
ON CONFLICT (user_id, role_id) DO UPDATE SET updated_at = excluded.updated_at;

-- GetUserRoles returns the roles for a user.
--
-- name: GetUserRoles :many
SELECT role_id
FROM user_roles
WHERE user_id = :user_id;

-- RemoveUserRole removes a role from a user.
--
-- name: RemoveUserRole :exec
DELETE
FROM user_roles
WHERE user_id = :user_id
  AND role_id = :role_id;
