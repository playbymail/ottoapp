-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- Games stores the data for an instance of a game.
CREATE TABLE games
(
    game_id     INTEGER PRIMARY KEY AUTOINCREMENT,

    code        TEXT    NOT NULL, -- 0300, 0301
    description TEXT    NOT NULL,

    active_turn TEXT    NOT NULL, -- YYYY-MM
    setup_turn  TEXT    NOT NULL, -- YYYY-MM
    orders_due  INTEGER NOT NULL, -- (unix seconds, UTC)

    -- audit (unix seconds, UTC)
    created_at  INTEGER NOT NULL, -- set in app
    updated_at  INTEGER NOT NULL, -- set in app

    UNIQUE (code)
);

CREATE TABLE game_turns
(
    game_id    INTEGER NOT NULL,
    turn       TEXT    NOT NULL, -- YYYY-MM

    turn_year  INTEGER NOT NULL DEFAULT 1 CHECK (turn_year > 0),
    turn_month INTEGER NOT NULL DEFAULT 1 CHECK (1 <= turn_month AND turn_month <= 12),
    turn_no    INTEGER NOT NULL DEFAULT 1 CHECK (turn_no >= 0),

    -- audit (unix seconds, UTC)
    created_at INTEGER NOT NULL, -- set in app
    updated_at INTEGER NOT NULL, -- set in app

    PRIMARY KEY (game_id, turn),
    UNIQUE (game_id, turn_year, turn_month),
    UNIQUE (game_id, turn_no),

    FOREIGN KEY (game_id)
        REFERENCES games (game_id)
        ON DELETE CASCADE
);
