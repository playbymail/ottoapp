-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- Clans stores the data for a clan. Clans are identified by number.
-- Don't confuse clans with tribes, which are elements controlled
-- by the clan.
--
-- A user may participate in many games, but may control only a single
-- clan in any game.
--
-- Clan number is unique within a game, but may be reused in other games.
--
-- The setup turn is the first turn in a game for a new player.
-- We track it because the report for that turn can be problematic.
-- It is often a truncated report containing the results of their setup
-- instructions and GM notes.
CREATE TABLE clans
(
    clan_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id    INTEGER NOT NULL, -- game that the clan participates in
    user_id    INTEGER NOT NULL, -- user that controls the clan

    clan       INTEGER NOT NULL CHECK (clan BETWEEN 1 and 999),
    is_active  BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    setup_turn TEXT    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL, -- set in app

    UNIQUE (user_id, game_id),
    UNIQUE (game_id, clan),
    FOREIGN KEY (game_id)
        REFERENCES games (game_id)
        ON DELETE CASCADE,
    FOREIGN KEY (game_id, setup_turn)
        REFERENCES game_turns (game_id, turn)
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
