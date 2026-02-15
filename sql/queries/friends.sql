-- name: CreateFriend :one
INSERT INTO friends (name)
VALUES ($1)
RETURNING *;

-- name: ListFriends :many
SELECT * FROM friends ORDER BY created_at DESC;

-- name: GetFriend :one
SELECT * FROM friends WHERE id = $1;

-- name: DeleteFriend :exec
DELETE FROM friends WHERE id = $1;
