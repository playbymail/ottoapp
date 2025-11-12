--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateUserSecret creates a secrets record for the user.
-- The password is stored as a bcrypt hash.
--
-- name: CreateUserSecret :exec
INSERT INTO user_secrets (user_id,
                          hashed_password,
                          last_login,
                          created_at,
                          updated_at)
VALUES (:user_id,
        :hashed_password,
        :last_login,
        :created_at,
        :updated_at);

-- GetUserRoles returns the roles for a user.
--
-- name: GetUserRoles :many
SELECT role_id
FROM user_roles
WHERE user_id = :user_id;

-- GetUserSecret returns the password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: GetUserSecret :one
SELECT hashed_password
FROM user_secrets
WHERE user_id = :user_id;

-- UpdateUserSecret updates password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: UpdateUserSecret :exec
UPDATE user_secrets
SET hashed_password = :hashed_password,
    updated_at      = :updated_at
WHERE user_id = :user_id;

-- AssignUserRole assigns a role to a user.
--
-- name: AssignUserRole :exec
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (:user_id, :role_id, :created_at, :updated_at);

-- RemoveUserRole removes a role from a user.
--
-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = :user_id AND role_id = :role_id;
