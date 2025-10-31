--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

CREATE TABLE schema_migrations
(
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    applied_at TEXT NOT NULL
);

CREATE TABLE schema_version
(
    id         INTEGER PRIMARY KEY,
    version    INTEGER NOT NULL UNIQUE,
    applied_at TEXT    NOT NULL
);

INSERT INTO schema_version (version, applied_at)
VALUES (1, current_timestamp);
