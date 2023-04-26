// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: attempt.sql

package db

import (
	"context"
)

const createAttempt = `-- name: CreateAttempt :one
INSERT INTO attempts (
  session_id,
  id,
  target_number,
  num_of_dice_throw,
  first_dice_throw_value,
  second_dice_throw_value
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id, session_id, target_number, num_of_dice_throw, first_dice_throw_value, second_dice_throw_value, created_at
`

type CreateAttemptParams struct {
	SessionID            int64  `json:"session_id"`
	ID                   string `json:"id"`
	TargetNumber         int16  `json:"target_number"`
	NumOfDiceThrow       int16  `json:"num_of_dice_throw"`
	FirstDiceThrowValue  int16  `json:"first_dice_throw_value"`
	SecondDiceThrowValue int16  `json:"second_dice_throw_value"`
}

func (q *Queries) CreateAttempt(ctx context.Context, arg CreateAttemptParams) (Attempt, error) {
	row := q.db.QueryRowContext(ctx, createAttempt,
		arg.SessionID,
		arg.ID,
		arg.TargetNumber,
		arg.NumOfDiceThrow,
		arg.FirstDiceThrowValue,
		arg.SecondDiceThrowValue,
	)
	var i Attempt
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.TargetNumber,
		&i.NumOfDiceThrow,
		&i.FirstDiceThrowValue,
		&i.SecondDiceThrowValue,
		&i.CreatedAt,
	)
	return i, err
}

const getAttempt = `-- name: GetAttempt :one
SELECT id, session_id, target_number, num_of_dice_throw, first_dice_throw_value, second_dice_throw_value, created_at FROM attempts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAttempt(ctx context.Context, id string) (Attempt, error) {
	row := q.db.QueryRowContext(ctx, getAttempt, id)
	var i Attempt
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.TargetNumber,
		&i.NumOfDiceThrow,
		&i.FirstDiceThrowValue,
		&i.SecondDiceThrowValue,
		&i.CreatedAt,
	)
	return i, err
}

const updateNumOfDiceThrown = `-- name: UpdateNumOfDiceThrown :one
UPDATE attempts SET num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING id, session_id, target_number, num_of_dice_throw, first_dice_throw_value, second_dice_throw_value, created_at
`

func (q *Queries) UpdateNumOfDiceThrown(ctx context.Context, id string) (Attempt, error) {
	row := q.db.QueryRowContext(ctx, updateNumOfDiceThrown, id)
	var i Attempt
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.TargetNumber,
		&i.NumOfDiceThrow,
		&i.FirstDiceThrowValue,
		&i.SecondDiceThrowValue,
		&i.CreatedAt,
	)
	return i, err
}

const updateValueofFirstDiceThrown = `-- name: UpdateValueofFirstDiceThrown :one
UPDATE attempts SET first_dice_throw_value = $2, num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING id, session_id, target_number, num_of_dice_throw, first_dice_throw_value, second_dice_throw_value, created_at
`

type UpdateValueofFirstDiceThrownParams struct {
	ID                  string `json:"id"`
	FirstDiceThrowValue int16  `json:"first_dice_throw_value"`
}

func (q *Queries) UpdateValueofFirstDiceThrown(ctx context.Context, arg UpdateValueofFirstDiceThrownParams) (Attempt, error) {
	row := q.db.QueryRowContext(ctx, updateValueofFirstDiceThrown, arg.ID, arg.FirstDiceThrowValue)
	var i Attempt
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.TargetNumber,
		&i.NumOfDiceThrow,
		&i.FirstDiceThrowValue,
		&i.SecondDiceThrowValue,
		&i.CreatedAt,
	)
	return i, err
}

const updateValueofSecondDiceThrown = `-- name: UpdateValueofSecondDiceThrown :one
UPDATE attempts SET second_dice_throw_value = $2, num_of_dice_throw = num_of_dice_throw + 1
WHERE id = $1
RETURNING id, session_id, target_number, num_of_dice_throw, first_dice_throw_value, second_dice_throw_value, created_at
`

type UpdateValueofSecondDiceThrownParams struct {
	ID                   string `json:"id"`
	SecondDiceThrowValue int16  `json:"second_dice_throw_value"`
}

func (q *Queries) UpdateValueofSecondDiceThrown(ctx context.Context, arg UpdateValueofSecondDiceThrownParams) (Attempt, error) {
	row := q.db.QueryRowContext(ctx, updateValueofSecondDiceThrown, arg.ID, arg.SecondDiceThrowValue)
	var i Attempt
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.TargetNumber,
		&i.NumOfDiceThrow,
		&i.FirstDiceThrowValue,
		&i.SecondDiceThrowValue,
		&i.CreatedAt,
	)
	return i, err
}