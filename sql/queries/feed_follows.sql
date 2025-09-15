-- name: CreateFeedFollows :one
WITH ins AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
        VALUES($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT ins.*, users.name AS username, feeds.name AS feedname FROM ins
    JOIN users ON users.id = ins.user_id
    JOIN feeds ON feeds.id = ins.feed_id;


