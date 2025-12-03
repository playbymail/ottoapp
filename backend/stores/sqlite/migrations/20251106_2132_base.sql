--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- The Users table stores user data.
-- Handle and Email are expected to be lower-cased.
CREATE TABLE users
(
    user_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    handle     TEXT    NOT NULL UNIQUE,
    username   TEXT    NOT NULL UNIQUE,
    email      TEXT    NOT NULL UNIQUE,
    timezone   TEXT    NOT NULL, -- IANA zone name

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL  -- set in app
);

-- The sysop is a required user for batch operations and system maintenance.
insert into users (user_id, username, handle, email, timezone, created_at, updated_at)
values (1, 'sysop', 'sysop', 'sysop', 'America/Panama', 0, 0);

-- The User_Secrets table stores credentials for authentication
CREATE TABLE user_secrets
(
    user_id            INTEGER PRIMARY KEY,
    hashed_password    TEXT    NOT NULL,
    plaintext_password TEXT,
    last_login         INTEGER NOT NULL, -- unix timestamp, must be UTC

    -- audit (unix seconds, UTC)
    created_at         INTEGER NOT NULL, -- set in app
    updated_at         INTEGER NOT NULL, -- set in app

    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);

-- The sysop is a required user for batch operations and system maintenance.
-- It is not allowed to log in, so we set the hashed password to an invalid value.
insert into user_secrets (user_id, hashed_password, plaintext_password, last_login, created_at, updated_at)
select user_id,
       '*',
       '*',
       0,
       users.created_at,
       users.updated_at
from users;

-- The Roles tables defines roles for authorization. The application
-- is responsible for implementing the "meaning" of each role.
CREATE TABLE roles
(
    role_id     TEXT PRIMARY KEY,
    is_active   BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    description TEXT    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at  INTEGER NOT NULL, -- set in app
    updated_at  INTEGER NOT NULL  -- set in app
);

-- define roles for the application
insert into roles (role_id, is_active, description, created_at, updated_at)
VALUES ('active', 1, 'active user role', 0, 0),
       ('sysop', 1, 'sysop role', 0, 0),
       ('admin', 1, 'administrator role', 0, 0),
       ('gm', 1, 'gm role', 0, 0),
       ('user', 1, 'user role', 0, 0),
       ('player', 1, 'player role', 0, 0),
       ('guest', 1, 'guest / anonymous visitor role', 0, 0),
       ('tn3', 1, 'game TN3 role', 0, 0),
       ('tn3.1', 1, 'game TN3.1 role', 0, 0)
;

-- The User_Roles table assigns roles to users.
CREATE TABLE user_roles
(
    user_id    INTEGER NOT NULL,
    role_id    TEXT    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL, -- set in app

    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE,
    FOREIGN KEY (role_id)
        REFERENCES roles (role_id)
        ON DELETE CASCADE
);

-- The sysop is a required user for batch operations and system maintenance.
insert into user_roles (user_id, role_id, created_at, updated_at)
select user_id, role_id, users.created_at, users.updated_at
from users
         cross join (select roles.role_id
                     from roles
                     where role_id in ('active', 'sysop'))
where users.username = 'sysop';

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

-- Games stores the data for an instance of a game.
--
-- The setup turn is the first turn in a new game.
CREATE TABLE games
(
    game_id          TEXT PRIMARY KEY,
    description      TEXT    NOT NULL UNIQUE,
    setup_turn_no    INTEGER NOT NULL CHECK (setup_turn_no >= 0),
    setup_turn_year  INTEGER NOT NULL CHECK (setup_turn_year BETWEEN 899 AND 9999),
    setup_turn_month INTEGER NOT NULL CHECK (setup_turn_month BETWEEN 1 AND 12),
    is_active        BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    -- audit (unix seconds, UTC)
    created_at       INTEGER NOT NULL, -- set in app
    updated_at       INTEGER NOT NULL  -- set in app
);

-- Clans stores the data for a clan. Clans are identified by number.
-- Don't confuse clans with tribes, which are elements controlled
-- by the clan.
--
-- A user may participate in many games, but may control only a single clan in any game.
--
-- Clan number is unique within a game, but may be reused in separate game.
--
-- The setup turn number is the first turn in a game for a new player.
-- We track it because the report for that turn is problematic.
-- It is a truncated report containing the results of their setup
-- instructions and GM notes.
CREATE TABLE clans
(
    clan_id       INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id       TEXT    NOT NULL, -- game that the clan participates in
    user_id       INTEGER NOT NULL, -- user that controls the clan

    clan          INTEGER NOT NULL CHECK (clan BETWEEN 1 and 999),
    setup_turn_no INTEGER NOT NULL CHECK (setup_turn_no >= 0),
    is_active     BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    -- audit (unix seconds, UTC)
    created_at    INTEGER NOT NULL, -- set in app
    updated_at    INTEGER NOT NULL, -- set in app

    UNIQUE (user_id, game_id),
    UNIQUE (game_id, clan),
    FOREIGN KEY (game_id)
        REFERENCES games (game_id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);

-- Elements stores meta-data for elements. Elements are the
-- type of unit that can be controlled by a Clan. There is
-- an Element that has a type of 'Element.'
CREATE TABLE elements
(
    element_type   TEXT PRIMARY KEY,
    element_suffix TEXT    NOT NULL UNIQUE,
    description    TEXT    NOT NULL UNIQUE,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL, -- set in app
    updated_at     INTEGER NOT NULL  -- set in app
);

INSERT INTO elements (element_type, element_suffix, description, created_at, updated_at)
VALUES ('TRIBE', '', 'Tribe element', 0, 0);
INSERT INTO elements (element_type, element_suffix, description, created_at, updated_at)
VALUES ('COURIER', 'c', 'Courier element', 0, 0);
INSERT INTO elements (element_type, element_suffix, description, created_at, updated_at)
VALUES ('ELEMENT', 'e', 'Element element', 0, 0);
INSERT INTO elements (element_type, element_suffix, description, created_at, updated_at)
VALUES ('FLEET', 'f', 'Fleet element', 0, 0);
INSERT INTO elements (element_type, element_suffix, description, created_at, updated_at)
VALUES ('GARRISON', 'g', 'Garrison element', 0, 0);
