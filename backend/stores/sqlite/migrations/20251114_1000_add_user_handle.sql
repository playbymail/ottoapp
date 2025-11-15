--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- Add handle column to users table
ALTER TABLE users ADD COLUMN handle TEXT;

-- Update existing users with their handles
UPDATE users SET handle = 'sysop' WHERE username = 'sysop';
UPDATE users SET handle = 'penguin' WHERE username = 'penguin';
UPDATE users SET handle = 'catbird' WHERE username = 'catbird';
UPDATE users SET handle = printf('h%06d', user_id) WHERE handle IS NULL;

-- Now make the column NOT NULL and add unique constraint
-- SQLite doesn't support ALTER COLUMN, so we need to recreate the constraint
-- by creating a new table and copying data
CREATE TABLE users_new
(
    user_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT UNIQUE NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    handle     TEXT UNIQUE NOT NULL,
    timezone   TEXT        NOT NULL,

    -- audit (unix seconds, UTC)
    created_at INTEGER     NOT NULL,
    updated_at INTEGER     NOT NULL
);

-- Copy data from old table to new table
INSERT INTO users_new (user_id, username, email, handle, timezone, created_at, updated_at)
SELECT user_id, username, email, handle, timezone, created_at, updated_at
FROM users;

-- Drop old table and rename new table
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;
