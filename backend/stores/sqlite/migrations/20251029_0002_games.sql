-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE game_turns
(
    turn_id    INTEGER PRIMARY KEY AUTOINCREMENT,

    code       TEXT    NOT NULL, -- YYYY-MM
    turn_year  INTEGER NOT NULL DEFAULT 1 CHECK (turn_year > 0),
    turn_month INTEGER NOT NULL DEFAULT 1 CHECK (1 <= turn_month AND turn_month <= 12),
    turn_no    INTEGER NOT NULL DEFAULT 1 CHECK (turn_no >= 0),

    orders_due INTEGER NOT NULL, -- (unix seconds, UTC)

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL  -- set in app
);

-- Games stores the data for an instance of a game.
CREATE TABLE games
(
    game_id        INTEGER PRIMARY KEY AUTOINCREMENT,

    code           TEXT    NOT NULL UNIQUE, -- 0300, 0301
    description    TEXT    NOT NULL UNIQUE,
    is_active      BOOL    NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    active_turn_id INTEGER NOT NULL,
    setup_turn_id  INTEGER NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL,        -- set in app
    updated_at     INTEGER NOT NULL,        -- set in app

    FOREIGN KEY (active_turn_id)
        REFERENCES game_turns (turn_id)
        ON DELETE CASCADE,
    FOREIGN KEY (setup_turn_id)
        REFERENCES game_turns (turn_id)
        ON DELETE CASCADE
);

