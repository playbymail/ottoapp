--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE users
(
    user_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    handle     TEXT UNIQUE NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    timezone   TEXT        NOT NULL,                           -- IANA zone name
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP  -- sqlite timestamp should be UTC
);

insert into users (handle, email, timezone)
values ('sysop', 'sysop', 'America/Panama');

INSERT INTO schema_version (version)
VALUES (2);
