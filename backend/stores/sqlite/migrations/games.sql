--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

CREATE TABLE games
(
    game_id     INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT UNIQUE NOT NULL,
    is_active   BOOL        NOT NULL DEFAULT 0,

    -- columns for auditing
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

insert into games (description, is_active)
values ('TN3', 1);

CREATE TABLE clans
(
    clan_id    INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER   NOT NULL,
    game_id    TEXT      NOT NULL,
    clan       TEXT      NOT NULL,

    -- columns for auditing
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (user_id, game_id),
    UNIQUE (game_id, clan),
    FOREIGN KEY (game_id) REFERENCES games (game_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

insert into clans (user_id, game_id, clan)
select users.user_id, games.game_id, '0000'
from users,
     games
where users.email = 'sysop'
  and games.is_active = 1;
