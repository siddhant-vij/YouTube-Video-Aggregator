-- name: InsertChannelFollow :exec
INSERT INTO channel_follows
  (id, created_at, updated_at, user_id, channel_id)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllChannelFollows :many
SELECT * FROM channel_follows;