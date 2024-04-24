-- name: InsertFeed :exec
INSERT INTO feeds
  (id, created_at, updated_at, name, url)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: UpdateLastFetchedAt :exec
UPDATE feeds
SET
  last_fetched_at = $2
WHERE id = $1;