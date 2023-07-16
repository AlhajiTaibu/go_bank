-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY created_at;

-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

-- name: UpdateTransfer :one
UPDATE transfers
  set from_account_id = $2,
  to_account_id = $3,
  amount = $4
WHERE id = $1
RETURNING *;