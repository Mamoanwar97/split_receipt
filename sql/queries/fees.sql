-- name: CreateFixedFee :one
INSERT INTO receipt_fixed_fees (receipt_id, name, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListFixedFeesByReceipt :many
SELECT * FROM receipt_fixed_fees WHERE receipt_id = $1 ORDER BY name;

-- name: DeleteFixedFee :exec
DELETE FROM receipt_fixed_fees WHERE id = $1;

-- name: CreatePercentageFee :one
INSERT INTO receipt_percentage_fees (receipt_id, name, percentage, cap_amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListPercentageFeesByReceipt :many
SELECT * FROM receipt_percentage_fees WHERE receipt_id = $1 ORDER BY name;

-- name: DeletePercentageFee :exec
DELETE FROM receipt_percentage_fees WHERE id = $1;
