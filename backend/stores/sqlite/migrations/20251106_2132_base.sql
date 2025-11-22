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

-- The Document_Contents table stores the data for a document.
-- One row per unique file (deduped by contents_hash).
--
-- MIME Types
--  DOCX        application/vnd.openxmlformats-officedocument.wordprocessingml.document
--  TURN_REPORT application/tn-3.0
--  TURN_REPORT application/tn-3.1
--  WXX         application/wxx.xml
--
-- Note: the game engine controls deletion of document contents.
-- They can be deleted when there are no child documents rows.
CREATE TABLE document_contents
(
    contents_hash  TEXT PRIMARY KEY,
    content_length INTEGER NOT NULL, -- size in bytes
    mime_type      TEXT    NOT NULL,
    contents       BLOB    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL, -- set in app
    updated_at     INTEGER NOT NULL  -- set in app
);

-- The Documents table contains meta-data for documents
-- (e.g., turn reports, maps).
--
-- The Document_Type column is used to categorize or filter
-- documents in the API.
--
-- Note: the upload process must untaint or reject the document_name.
CREATE TABLE documents
(
    document_id   INTEGER PRIMARY KEY AUTOINCREMENT,
    clan_id       INTEGER NOT NULL, -- clan that owns this document

    -- permissions on the document. this is supposed to let
    -- a user share a document with an ally or the system to
    -- create a shared read-only document (e.g., the map key).
    can_read      BOOL    NOT NULL DEFAULT 0 CHECK (can_read IN (0, 1)),
    can_write     BOOL    NOT NULL DEFAULT 0 CHECK (can_write IN (0, 1)),
    can_delete    BOOL    NOT NULL DEFAULT 0 CHECK (can_delete IN (0, 1)),
    can_share     BOOL    NOT NULL DEFAULT 0 CHECK (can_share IN (0, 1)),

    document_name TEXT    NOT NULL, -- user's name for the doc
    document_type TEXT    NOT NULL, -- report, map, etc

    contents_hash TEXT    NOT NULL, -- 64-char sha256 hex

    -- audit (unix seconds, UTC)
    created_at    INTEGER NOT NULL, -- set in app
    updated_at    INTEGER NOT NULL, -- set in app

    -- prevent users from uploading the same document multiple times
    UNIQUE (clan_id, contents_hash),

    FOREIGN KEY (clan_id)
        REFERENCES clans (clan_id),
    FOREIGN KEY (contents_hash)
        REFERENCES document_contents (contents_hash)
        ON DELETE CASCADE
);

-- index for "show me all docs I own"
CREATE INDEX idx_documents_owner
    ON documents (clan_id);

-- The Document_Shares table is a bridge table for sharing documents.
CREATE TABLE document_shares
(
    document_id INTEGER PRIMARY KEY,
    clan_id     INTEGER NOT NULL,

    can_read    BOOL    NOT NULL DEFAULT 0 CHECK (can_read IN (0, 1)),
    can_delete  BOOL    NOT NULL DEFAULT 0 CHECK (can_delete IN (0, 1)),

    -- audit (unix seconds, UTC)
    created_at  INTEGER NOT NULL, -- set in app
    updated_at  INTEGER NOT NULL, -- set in app

    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE,
    FOREIGN KEY (clan_id)
        REFERENCES clans (clan_id)
);

-- index for "show me all docs shared with me"
CREATE INDEX idx_documents_shared
    ON document_shares (clan_id);

-- Clan_Documents_VW returns all documents owned by or shared with a clan.
CREATE VIEW clan_documents_vw
            (
             document_id,
             clan_id,
             can_read,
             can_write,
             can_delete,
             can_share,
             document_name,
             document_type,
             contents_hash,
             owner_id,
             is_shared,
             created_at,
             updated_at
                )
AS
SELECT d.document_id,
       d.clan_id,
       CASE WHEN s.clan_id IS NULL THEN d.can_read ELSE s.can_read END     AS can_read,
       CASE WHEN s.clan_id IS NULL THEN d.can_write ELSE 0 END             AS can_write,
       CASE WHEN s.clan_id IS NULL THEN d.can_delete ELSE s.can_delete END AS can_delete,
       CASE WHEN s.clan_id IS NULL THEN d.can_share ELSE 0 END             AS can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       CASE WHEN s.clan_id IS NULL THEN d.clan_id ELSE s.clan_id END       AS owner_id,
       CASE WHEN s.clan_id IS NULL THEN 0 ELSE 1 END                       AS is_shared,
       d.created_at,
       d.updated_at
FROM documents AS d
         LEFT JOIN document_shares AS s
                   ON d.document_id = s.document_id;

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
