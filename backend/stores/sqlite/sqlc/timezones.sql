--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- GetActiveTimezones returns the timezones that are used by users.
--
-- name: GetActiveTimezones :many
SELECT distinct timezone
FROM users
ORDER BY 1;
