-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE schema_migrations
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    migration_id TEXT    NOT NULL,
    file_name    TEXT    NOT NULL UNIQUE,
    applied_at   INTEGER NOT NULL, -- unix seconds, UTC

    -- audit (unix seconds, UTC)
    created_at   INTEGER NOT NULL, -- set in app
    updated_at   INTEGER NOT NULL, -- set in app

    UNIQUE (migration_id)
);

CREATE TABLE config
(
    key        TEXT    NOT NULL,
    value      TEXT    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL, -- set in app

    PRIMARY KEY (key)
);

INSERT INTO config (key, value, created_at, updated_at)
VALUES ('schema.version', '20251029_0000', 0, 0);
