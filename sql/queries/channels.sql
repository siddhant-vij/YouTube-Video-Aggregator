-- name: InsertChannel :exec
INSERT INTO channels
  (id, created_at, updated_at, name, url, last_fetched_at)
VALUES
  ($1, $2, $3, $4, $5, $6);

-- name: GetAllChannels :many
SELECT * FROM channels;

-- name: GetNumChannelsByCreatedAt :many
SELECT * FROM channels
ORDER BY created_at DESC
LIMIT $1;

-- name: UpdateLastFetchedAt :exec
UPDATE channels
SET
  last_fetched_at = $2,
  updated_at = $2
WHERE id = $1;