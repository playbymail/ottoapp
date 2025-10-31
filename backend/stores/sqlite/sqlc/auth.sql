--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateUserSecrets creates a secrets record for the user.
-- The password is stored as a bcrypt hash.
--
-- name: CreateUserSecrets :exec
INSERT INTO user_secrets (user_id,
                   hashed_password,
                   last_login)
VALUES (:user_id,
        :hashed_password,
        :last_login);

-- GetUserPassword returns the password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: GetUserSecrets :one
SELECT hashed_password
FROM user_secrets
WHERE user_id = :user_id;

-- UpdateUserPassword updates password for a user.
-- The password is stored as a bcrypt hash.
--
-- name: UpdateUserPassword :exec
UPDATE user_secrets
SET hashed_password = :hashed_password,
    updated_at      = CURRENT_TIMESTAMP
WHERE user_id = :user_id;
