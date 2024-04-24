-- name: InsertFeedFollow :exec
INSERT INTO feed_follows
  (id, created_at, updated_at, user_id, feed_id)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllFeedFollows :many
SELECT * FROM feed_follows;