--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- name: GetDatabaseVersion :one
SELECT value
FROM config
WHERE key = 'schema.version';

-- name: GetDatabaseMigrationsApplied :many
SELECT id, migration_id, file_name, applied_at
FROM schema_migrations;