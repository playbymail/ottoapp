-- CreateUserGameClan has two business rules
--  user can have at most one clan per game.
--  clan number can be used by at most one user per game.
-- The upsert key is "user_id, game_id," preventing a user
-- from claiming multiple clans in a game. If a user tries
-- claiming an existing clan in a game, it will fail, not
-- silently clobber another user's.
--
-- name: CreateGameUserClan :one
INSERT INTO clans (game_id,
                   user_id,
                   clan,
                   setup_turn,
                   created_at,
                   updated_at)
VALUES (:game_id,
        :user_id,
        :clan,
        :setup_turn,
        :created_at,
        :updated_at)
ON CONFLICT (game_id, user_id) DO UPDATE
    SET clan       = excluded.clan,
        setup_turn = excluded.setup_turn,
        updated_at = excluded.updated_at
RETURNING clan_id;

-- name: GetClan :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn,
       is_active
FROM clans
WHERE clan_id = :clan_id;

-- name: GetClanByGameUser :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn,
       is_active
FROM clans
WHERE game_id = :game_id
  AND user_id = :user_id;

-- name: GetClanByGameClanNo :one
SELECT game_id,
       user_id,
       clan_id,
       clan,
       setup_turn,
       is_active
FROM clans
WHERE game_id = :game_id
  AND clan = :clan_no;

-- name: ReadClanByGameIdClanNo :one
SELECT game_id, user_id, clan_id, clan
FROM clans
WHERE game_id = :game_id
  AND clan = :clan_no;

-- name: ReadClansByGame :many
SELECT game_id, user_id, clan_id, clan
FROM clans
WHERE game_id = :game_id
  AND is_active = 1
ORDER BY clans.clan;

-- name: RemoveClan :exec
DELETE
FROM clans
WHERE clan_id = :clan_id;
