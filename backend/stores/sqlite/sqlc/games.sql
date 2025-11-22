--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- name: UpsertGame :one
INSERT INTO games (game_id,
                   description,
                   setup_turn_no,
                   setup_turn_year,
                   setup_turn_month,
                   is_active,
                   created_at,
                   updated_at)
VALUES (:game_id,
        :description,
        :setup_turn_no,
        :setup_turn_year,
        :setup_turn_month,
        :is_active,
        :created_at,
        :updated_at)
ON CONFLICT (game_id) DO UPDATE
    SET description      = excluded.description,
        setup_turn_no    = excluded.setup_turn_no,
        setup_turn_year  = excluded.setup_turn_year,
        setup_turn_month = excluded.setup_turn_month,
        is_active        = excluded.is_active,
        updated_at       = excluded.updated_at
RETURNING game_id;

-- name: GetGame :one
SELECT game_id,
       description,
       setup_turn_year,
       setup_turn_month,
       setup_turn_no,
       is_active
FROM games
WHERE game_id = :game_id;


-- UpsertGameUserClan has two business rules
--  user can have at most one clan per game.
--  clan number can be used by at most one user per game.
-- The upsert key is "user_id, game_id," preventing a user
-- from claiming multiple clans in a game. If a user tries
-- claiming an existing clan in a game, it will fail, not
-- silently clobber another user's.
--
-- name: UpsertGameUserClan :one
INSERT INTO clans (game_id, user_id, clan, setup_turn_no, created_at, updated_at)
VALUES (:game_id, :user_id, :clan, :setup_turn_no, :created_at, :updated_at)
ON CONFLICT (user_id, game_id) DO UPDATE SET clan          = excluded.clan,
                                             setup_turn_no = excluded.setup_turn_no,
                                             updated_at    = excluded.updated_at
RETURNING clan_id;

-- name: GetClan :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn_no,
       is_active
FROM clans
WHERE clan_id = :clan_id;

-- name: GetClanByGameUser :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn_no,
       is_active
FROM clans
WHERE game_id = :game_id
  AND user_id = :user_id;

-- name: GetClanByGameClanNo :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn_no,
       is_active
FROM clans
WHERE game_id = :game_id
  AND clan = :clan_no;

-- name: RemoveClan :exec
DELETE
FROM clans
WHERE clan_id = :clan_id;
