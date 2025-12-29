--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- name: CreateConfigKeyValue :exec
INSERT INTO config (key, value, created_at, updated_at)
VALUES (:key, :value, :created_at, :updated_at);

-- name: ReadConfigKeyValue :one
SELECT value
FROM config
WHERE key = :key;

-- name: UpdateConfigKeyValue :exec
UPDATE config
SET value      = :value,
    updated_at = :updated_at
WHERE key = :key;

-- name: UpsertConfigKeyValue :one
INSERT INTO config (key, value, created_at, updated_at)
VALUES (:key, :value, :created_at, :updated_at)
ON CONFLICT (key)
    DO UPDATE
    SET value      = excluded.value,
        updated_at = excluded.updated_at
RETURNING key;

-- name: DeleteConfigKeyValue :exec
DELETE
FROM config
WHERE key = :key;
