-- name: CreateSession :one
INSERT INTO sessions (
  is_active, owner, current_attempt_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: UpdateSessionCurrentAttemptId :one
UPDATE sessions SET current_attempt_id = $2
WHERE id = $1
RETURNING *;

-- name: EndUserGameSession :one
UPDATE sessions SET is_active = FALSE
WHERE id = $1 AND is_active = TRUE
RETURNING *;

-- name: GetUserWithSession :one
SELECT  s.owner, s.current_attempt_id, s.is_active, s.id, a.target_number, a.num_of_dice_throw, a.first_dice_throw_value, a.second_dice_throw_value FROM sessions s
INNER JOIN attempts a
ON s.current_attempt_id = a.id
WHERE owner = $1 
AND s.is_active = TRUE
ORDER BY s.created_at
LIMIT 1;

