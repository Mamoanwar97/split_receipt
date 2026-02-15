-- name: CreateMeal :one
INSERT INTO meals (receipt_id, name, quantity, total_price, price_per_unit)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListMealsByReceipt :many
SELECT * FROM meals WHERE receipt_id = $1 ORDER BY created_at DESC;

-- name: GetMeal :one
SELECT * FROM meals WHERE id = $1;

-- name: DeleteMeal :exec
DELETE FROM meals WHERE id = $1;

-- name: AddFriendToMeal :exec
INSERT INTO meal_friends (meal_id, friend_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveFriendFromMeal :exec
DELETE FROM meal_friends WHERE meal_id = $1 AND friend_id = $2;

-- name: ListFriendsByMeal :many
SELECT f.* FROM friends f
JOIN meal_friends mf ON mf.friend_id = f.id
WHERE mf.meal_id = $1
ORDER BY f.name;

-- name: CountFriendsByMeal :one
SELECT count(*) FROM meal_friends WHERE meal_id = $1;

-- name: CountUniqueFriendsByReceipt :one
SELECT count(DISTINCT mf.friend_id)
FROM meal_friends mf
JOIN meals m ON m.id = mf.meal_id
WHERE m.receipt_id = $1;

-- name: ListMealsByReceiptAndFriend :many
SELECT m.* FROM meals m
JOIN meal_friends mf ON mf.meal_id = m.id
WHERE m.receipt_id = $1 AND mf.friend_id = $2;
