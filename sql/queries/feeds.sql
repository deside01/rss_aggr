-- name: CreateFeed :one
INSERT INTO feeds (
    id, user_id, url, name, created_at, updated_at
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetUserFeeds :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetNextFeeds :many
SELECT * FROM feeds
ORDER BY created_at ASC
LIMIT $1;

-- name: DeleteFeed :one
DELETE FROM feeds WHERE id = $1
RETURNING id;