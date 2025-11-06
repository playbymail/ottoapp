--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- GetConfigKeyValue does
--
-- name: GetConfigKeyValue :one
SELECT value
FROM config
WHERE key = :key;
