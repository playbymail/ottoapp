-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- The Users table stores user data.
-- Handle and Email are expected to be lower-cased.
CREATE TABLE users
(
    user_id            INTEGER PRIMARY KEY AUTOINCREMENT,
    handle             TEXT    NOT NULL UNIQUE,
    username           TEXT    NOT NULL UNIQUE,
    email              TEXT    NOT NULL UNIQUE,
    email_opt_in       BOOL    NOT NULL DEFAULT 0 CHECK (email_opt_in IN (0, 1)),
    timezone           TEXT    NOT NULL,           -- IANA zone name
    is_active          BOOL    NOT NULL DEFAULT 0 CHECK (is_active IN (0, 1)),
    is_admin           BOOL    NOT NULL DEFAULT 0 CHECK (is_admin IN (0, 1)),
    is_gm              BOOL    NOT NULL DEFAULT 0 CHECK (is_gm IN (0, 1)),
    is_guest           BOOL    NOT NULL DEFAULT 0 CHECK (is_guest IN (0, 1)),
    is_player          BOOL    NOT NULL DEFAULT 0 CHECK (is_player IN (0, 1)),
    is_service         BOOL    NOT NULL DEFAULT 0 CHECK (is_service IN (0, 1)),
    is_sysop           BOOL    NOT NULL DEFAULT 0 CHECK (is_sysop IN (0, 1)),
    is_user            BOOL    NOT NULL DEFAULT 0 CHECK (is_user IN (0, 1)),
    hashed_password    TEXT    NOT NULL DEFAULT '*',
    plaintext_password TEXT    NOT NULL DEFAULT '',
    last_login         INTEGER NOT NULL DEFAULT 0, -- unix timestamp, must be UTC

    -- audit (unix seconds, UTC)
    created_at         INTEGER NOT NULL,           -- set in app
    updated_at         INTEGER NOT NULL            -- set in app
);

-- The sysop is a required user for batch operations and system maintenance.
insert into users (user_id,
                   username,
                   handle,
                   email,
                   timezone,
                   is_active,
                   is_sysop,
                   created_at,
                   updated_at)
values (1,
        'sysop',
        'sysop',
        'sysop',
        'America/Panama',
        1,
        1,
        0,
        0);

-- The Sessions table holds data for each session.
CREATE TABLE sessions
(
    session_id TEXT PRIMARY KEY,
    csrf       TEXT    NOT NULL,
    user_id    INTEGER NOT NULL,
    expires_at INTEGER NOT NULL, -- unix seconds, UTC

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL, -- set in app

    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);

