-- name: CreateUser :one
INSERT INTO users (
  username
) VALUES (
  $1
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserTransactionLogs :many
SELECT  u.username, e.amount, w.asset, e.created_at as transaction_date, e.id,
       CASE
           WHEN e.amount > 0 THEN 'credit'
           WHEN e.amount < 0 THEN 'debit'
       END transaction_type
FROM users u
INNER JOIN wallets w
ON u.username = w.owner
INNER JOIN entries e
ON w.id = e.wallet_id
WHERE username = $1 
ORDER BY e.created_at;