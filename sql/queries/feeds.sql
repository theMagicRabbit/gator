-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeed :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds SET updated_at = $1, last_fetched_at = $1
WHERE id = $2
RETURNING *;

