-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY owner LIMIT $1 OFFSET $2;

-- name: CreateAccount :one
INSERT INTO accounts (
  owner, currency, balance
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE accounts
  set owner = $2,
  currency = $3,
  balance = $4
WHERE id = $1
RETURNING *;