-- name: InsertChannel :exec
INSERT INTO channels
  (id, created_at, updated_at, name, url)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllChannels :many
SELECT * FROM channels;

-- name: UpdateLastFetchedAt :exec
UPDATE channels
SET
  last_fetched_at = $2
WHERE id = $1;