--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- The Users table stores user data.
CREATE TABLE users
(
    user_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    handle     TEXT UNIQUE NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    timezone   TEXT        NOT NULL, -- IANA zone name

    -- audit (unix seconds, UTC)
    created_at INTEGER     NOT NULL, -- set in app
    updated_at INTEGER     NOT NULL  -- set in app
);

-- The sysop is a required user for batch operations and system maintenance.
insert into users (handle, email, timezone, created_at, updated_at)
values ('sysop', 'sysop', 'America/Panama', 0, 0);

-- The User_Secrets table stores credentials for authentication
CREATE TABLE user_secrets
(
    user_id         INTEGER PRIMARY KEY,
    hashed_password TEXT    NOT NULL,
    last_login      INTEGER NOT NULL, -- unix timestamp, must be UTC

    -- audit (unix seconds, UTC)
    created_at      INTEGER NOT NULL, -- set in app
    updated_at      INTEGER NOT NULL, -- set in app

    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);

-- The sysop is a required user for batch operations and system maintenance.
insert into user_secrets (user_id, hashed_password, last_login, created_at, updated_at)
select user_id,
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
where users.handle = 'sysop';

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
    game_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    description      TEXT    NOT NULL,
    setup_turn_no    INTEGER NOT NULL CHECK (setup_turn_no >= 0),
    setup_turn_year  INTEGER NOT NULL CHECK (setup_turn_year BETWEEN 899 AND 9999),
    setup_turn_month INTEGER NOT NULL CHECK (setup_turn_month BETWEEN 1 AND 12),
    is_active        BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    -- audit (unix seconds, UTC)
    created_at       INTEGER NOT NULL, -- set in app
    updated_at       INTEGER NOT NULL, -- set in app

    UNIQUE (description)
);

-- insert into games (game_id, description, setup_turn_no, setup_turn_year, setup_turn_month, is_active)
-- values ('0300', 'TN3', 0, 899, 12, 1),
--        ('0301', 'TN3.1', 0, 899, 12, 1);

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
    element_type   TEXT    NOT NULL,
    element_suffix TEXT    NOT NULL,
    description    TEXT    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL, -- set in app
    updated_at     INTEGER NOT NULL, -- set in app

    PRIMARY KEY (element_type),
    UNIQUE (element_suffix),
    UNIQUE (description)
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

-- The Documents table contains meta-data for documents (e.g., turn reports, maps).
-- One row per unique file (deduped by contents_hash).
--
-- MIME Types
--  DOCX        application/vnd.openxmlformats-officedocument.wordprocessingml.document
--  TURN_REPORT application/tn-3.0
--  TURN_REPORT application/tn-3.1
--  WXX         application/wxx.xml
CREATE TABLE documents
(
    document_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    mime_type      TEXT    NOT NULL,
    contents_hash  TEXT    NOT NULL UNIQUE, -- 64-char sha256 hex
    content_length INTEGER NOT NULL,        -- size in bytes

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL,        -- set in app
    updated_at     INTEGER NOT NULL         -- set in app
);

-- The Document_Contents table stores the data for a document.
-- Actual bytes, split out so we only store them once per document.
CREATE TABLE document_contents
(
    document_id INTEGER PRIMARY KEY,
    contents    BLOB    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at  INTEGER NOT NULL, -- set in app
    updated_at  INTEGER NOT NULL, -- set in app

    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE
);

-- The Document_ACL table tracks per-user ACL/view of the document.
-- ownership and other permissions on the document.
--
-- Each user who uploads or is granted access gets a row here.
--
-- Note: the game engine control deletion of documents. When there are no
-- Document_ACL entries for a document, it can be deleted.
CREATE TABLE document_acl
(
    document_id   INTEGER NOT NULL,
    user_id       INTEGER NOT NULL,

    document_name TEXT    NOT NULL, -- tainted: user's name for the doc
    created_by    INTEGER NOT NULL, -- who granted/created this ACL row

    is_owner      BOOL    NOT NULL DEFAULT 0 CHECK (is_owner IN (0, 1)),
    can_read      BOOL    NOT NULL DEFAULT 0 CHECK (can_read IN (0, 1)),
    can_write     BOOL    NOT NULL DEFAULT 0 CHECK (can_write IN (0, 1)),
    can_delete    BOOL    NOT NULL DEFAULT 0 CHECK (can_delete IN (0, 1)),

    -- audit (unix seconds, UTC)
    created_at    INTEGER NOT NULL, -- set in app
    updated_at    INTEGER NOT NULL, -- set in app

    PRIMARY KEY (document_id, user_id),
    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
);

-- index for "show me all docs I can see"
CREATE INDEX idx_document_acl_user ON document_acl (user_id);


-- The Turn_Reports table contains meta-data for turn reports.
--
-- We impose a constraint on the report - every section in it
-- must be for the same turn and all elements must be in the
-- same Clan.
CREATE TABLE turn_reports
(
    turn_report_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id        INTEGER NOT NULL,                      -- game the report was created for
    user_id        INTEGER NOT NULL,                      -- user that owns the report
    document_id    INTEGER NOT NULL,
    turn_no        INTEGER NOT NULL CHECK (turn_no >= 0), -- turn number from the report
    clan_id        INTEGER NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL,                      -- set in app
    updated_at     INTEGER NOT NULL,                      -- set in app

    FOREIGN KEY (clan_id)
        REFERENCES clans (clan_id)
        ON DELETE CASCADE,
    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE,
    FOREIGN KEY (game_id)
        REFERENCES games (game_id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);
