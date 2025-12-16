-- name: CreateGame :one
INSERT INTO games (code, description, active_turn_id, setup_turn_id, created_at, updated_at)
VALUES (:code, :description, :active_turn_id, :setup_turn_id, :created_at, :updated_at)
RETURNING game_id;

-- name: ReadGame :one
SELECT games.game_id,
       games.code,
       games.is_active,
       games.description,
       active_turn.turn_id    AS active_turn_id,
       active_turn.turn_year  AS active_turn_year,
       active_turn.turn_month AS active_turn_month,
       active_turn.turn_no    AS active_turn_no,
       active_turn.orders_due AS orders_due
FROM games,
     game_turns AS active_turn
WHERE games.game_id = :game_id
  AND active_turn.turn_id = games.active_turn;

-- name: ReadGames :many
SELECT games.game_id,
       games.code,
       games.is_active,
       games.description,
       active_turn.turn_id    AS active_turn_id,
       active_turn.turn_year  AS active_turn_year,
       active_turn.turn_month AS active_turn_month,
       active_turn.turn_no    AS active_turn_no,
       active_turn.orders_due AS orders_due
FROM games,
     game_turns AS active_turn
WHERE active_turn.turn_id = games.active_turn;

-- name: UpdateGameActiveTurn :exec
UPDATE games
SET active_turn_id = :active_turn_id,
    updated_at     = :updated_at
WHERE game_id = :game_id;

-- name: UpsertGame :one
INSERT INTO games (code, description, active_turn_id, setup_turn_id, created_at, updated_at)
VALUES (:code, :description, :active_turn_id, :setup_turn_id, :created_at, :updated_at)
ON CONFLICT (code) DO UPDATE
    SET description    = excluded.description,
        is_active      = excluded.is_active,
        active_turn_id = excluded.active_turn_id,
        setup_turn_id  = excluded.setup_turn_id,
        updated_at     = excluded.updated_at
RETURNING game_id;


-- name: DeleteGame :exec
DELETE
FROM games
WHERE game_id = :game_id;

