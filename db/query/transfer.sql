-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1;

-- name: ListOutTransfer :many
SELECT *
FROM transfers
WHERE from_account_id = $1
LIMIT $2 OFFSET $3;

-- name: ListInTransfer :many
SELECT *
FROM transfers
WHERE to_account_id = $1
LIMIT $2 OFFSET $3;

-- name: ListUnitedTransfer :many
SELECT *
FROM transfers
WHERE from_account_id = $1
  AND to_account_id = $2
LIMIT $3 OFFSET $4;

-- name: UpdateTransfer :one
UPDATE transfers
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE
FROM transfers
WHERE id = $1;