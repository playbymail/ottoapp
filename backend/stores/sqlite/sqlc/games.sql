-- name: CreateGame :one
INSERT INTO games (code, description, active_turn, setup_turn, orders_due, created_at, updated_at)
VALUES (:code, :description, :active_turn, :setup_turn, :orders_due, :created_at, :updated_at)
ON CONFLICT (code)
    DO UPDATE SET description = excluded.description,
                  active_turn = excluded.active_turn,
                  setup_turn  = excluded.setup_turn,
                  orders_due  = excluded.orders_due,
                  updated_at  = excluded.updated_at
RETURNING game_id;

-- name: CreateGameTurn :exec
INSERT INTO game_turns(game_id, turn, turn_year, turn_month, turn_no, created_at, updated_at)
VALUES (:game_id, :turn, :turn_year, :turn_month, :turn_no, :created_at, :updated_at)
ON CONFLICT (game_id, turn)
    DO UPDATE SET turn_year  = excluded.turn_year,
                  turn_month = excluded.turn,
                  turn_no    = excluded.turn_no,
                  updated_at = excluded.updated_at;;

-- name: ReadGame :one
SELECT games.game_id,
       games.code,
       games.description,
       games.active_turn,
       game_turns.turn_year,
       game_turns.turn_month,
       games.setup_turn,
       games.orders_due
FROM games,
     game_turns
WHERE games.game_id = :game_id
  AND games.game_id = game_turns.game_id
  AND game_turns.turn = games.active_turn;

-- name: ReadGames :many
SELECT games.game_id,
       games.code,
       games.description,
       games.active_turn,
       game_turns.turn_year,
       game_turns.turn_month,
       games.setup_turn,
       games.orders_due
FROM games,
     game_turns
WHERE games.game_id = game_turns.game_id
  AND game_turns.turn = games.active_turn;

-- name: UpdateGame :exec
UPDATE games
SET active_turn = :active_turn,
    setup_turn  = :setup_turn,
    orders_due  = :orders_due,
    updated_at  = :updated_at
WHERE game_id = :game_id;

-- name: UpdateGameActiveTurn :exec
UPDATE games
SET active_turn = :turn,
    updated_at  = :updated_at
WHERE game_id = :game_id;

-- name: DeleteGame :exec
DELETE
FROM games
WHERE game_id = :game_id;

