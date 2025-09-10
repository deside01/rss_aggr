-- name: FollowFeed :one
INSERT INTO feed_follows (
    id, user_id, feed_id, created_at, updated_at
)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserFollowedFeeds :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: DeleteFollowedFeed :one
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2
RETURNING id;