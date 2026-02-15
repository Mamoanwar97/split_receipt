-- name: CreateReceipt :one
INSERT INTO receipts (restaurant_name, date_time, total)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListReceipts :many
SELECT * FROM receipts ORDER BY created_at DESC;

-- name: GetReceipt :one
SELECT * FROM receipts WHERE id = $1;

-- name: DeleteReceipt :exec
DELETE FROM receipts WHERE id = $1;
