--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- Games stores the data for an instance of a game.
--
-- The setup turn is the first turn in a new game.
CREATE TABLE games
(
    game_id          INTEGER PRIMARY KEY AUTOINCREMENT,
    description      TEXT      NOT NULL,
    setup_turn_no    INTEGER   NOT NULL CHECK (setup_turn_no >= 0),
    setup_turn_year  INTEGER   NOT NULL CHECK (setup_turn_year BETWEEN 899 AND 9999),
    setup_turn_month INTEGER   NOT NULL CHECK (setup_turn_month BETWEEN 1 AND 12),
    is_active        BOOL      NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    -- columns for auditing
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    UNIQUE (description)
);

insert into games (game_id, description, setup_turn_no, setup_turn_year, setup_turn_month, is_active)
values ('0300', 'TN3', 0, 899, 12, 1),
       ('0301', 'TN3.1', 0, 899, 12, 1);

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
    game_id       TEXT      NOT NULL,                           -- game that the clan participates in
    user_id       INTEGER   NOT NULL,                           -- user that controls the clan

    clan          INTEGER   NOT NULL CHECK (clan BETWEEN 1 and 999),
    setup_turn_no INTEGER   NOT NULL CHECK (setup_turn_no >= 0),
    is_active     BOOL      NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    -- columns for auditing
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    UNIQUE (user_id, game_id),
    UNIQUE (game_id, clan),
    FOREIGN KEY (game_id) REFERENCES games (game_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- Elements stores meta-data for elements. Elements are the
-- type of unit that can be controlled by a Clan. There is
-- an Element that has a type of 'Element.'
CREATE TABLE elements
(
    element_type   TEXT      NOT NULL,
    element_suffix TEXT      NOT NULL,
    description    TEXT      NOT NULL,

    -- columns for auditing
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    PRIMARY KEY (element_type),
    UNIQUE (element_suffix),
    UNIQUE (description)
);

INSERT INTO elements (element_type, element_suffix, description)
VALUES ('TRIBE', '', 'Tribe element');
INSERT INTO elements (element_type, element_suffix, description)
VALUES ('COURIER', 'c', 'Courier element');
INSERT INTO elements (element_type, element_suffix, description)
VALUES ('ELEMENT', 'e', 'Element element');
INSERT INTO elements (element_type, element_suffix, description)
VALUES ('FLEET', 'f', 'Fleet element');
INSERT INTO elements (element_type, element_suffix, description)
VALUES ('GARRISON', 'g', 'Garrison element');

-- Documents contains meta-data for documents (e.g., turn reports) that are uploaded or created on the server.
CREATE TABLE documents
(
    document_id         TEXT      NOT NULL,                           -- SHA-256 hash of the document contents
    document_created_by INTEGER   NOT NULL,                           -- user that created the document
    document_created_at INTEGER   NOT NULL,                           -- unix timestamp, must always be UTC
    document_path       TEXT      NOT NULL,                           -- location of the document on the server

    -- columns for auditing
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    PRIMARY KEY (document_id),

    -- game engine must control deletion of documents; don't do it automatically
    FOREIGN KEY (document_created_by) REFERENCES users (user_id)
);

-- Turn_Reports contains meta-data for turn reports.
--
-- We impose a constraint on the report - every section in it must be for the same turn.
CREATE TABLE turn_reports
(
    turn_report_id INTEGER PRIMARY KEY AUTOINCREMENT,

    document_id    INTEGER NOT NULL,

    game_id        INTEGER NOT NULL,                      -- game the report was created for
    user_id        INTEGER NOT NULL,                      -- user that owns the report
    turn_no        INTEGER NOT NULL CHECK (turn_no >= 0), -- turn number from the report

    FOREIGN KEY (document_id) REFERENCES documents (document_id),
    FOREIGN KEY (game_id) REFERENCES games (game_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

-- Turn_Report_ACL controls which users have access to a turn report.
-- The user can either own the report or be a guest.
CREATE TABLE turn_report_acl
(
    turn_report_id INTEGER   NOT NULL,
    user_id        INTEGER   NOT NULL,
    is_owner       BOOL      NOT NULL DEFAULT 0 CHECK (is_owner IN (0, 1)),
    can_read       BOOL      NOT NULL DEFAULT 0 CHECK (can_read IN (0, 1)),
    can_write      BOOL      NOT NULL DEFAULT 0 CHECK (can_write IN (0, 1)),
    can_delete     BOOL      NOT NULL DEFAULT 0 CHECK (can_delete IN (0, 1)),

    -- columns for auditing
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    PRIMARY KEY (turn_report_id, user_id),
    FOREIGN KEY (turn_report_id) REFERENCES turn_reports (turn_report_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

INSERT INTO schema_version (version)
VALUES (6);

UPDATE config
SET VALUE = '20251105_0844',
    updated_at = CURRENT_TIMESTAMP
WHERE key = 'schema_version';
