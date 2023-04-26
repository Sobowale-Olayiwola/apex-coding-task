-- name: CreateAttempt :one
INSERT INTO attempts (
  session_id,
  id,
  target_number,
  num_of_dice_throw,
  first_dice_throw_value,
  second_dice_throw_value
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAttempt :one
SELECT * FROM attempts
WHERE id = $1 LIMIT 1;

-- name: UpdateNumOfDiceThrown :one
UPDATE attempts SET num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING *;

-- name: UpdateValueofFirstDiceThrown :one
UPDATE attempts SET first_dice_throw_value = $2, num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING *;

-- name: UpdateValueofSecondDiceThrown :one
UPDATE attempts SET second_dice_throw_value = $2, num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING *;