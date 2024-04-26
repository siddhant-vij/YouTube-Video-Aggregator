-- name: InsertChannelFollow :exec
INSERT INTO channel_follows
  (id, created_at, updated_at, user_id, channel_id)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllChannelFollows :many
SELECT * FROM channel_follows;

-- name: GetChannelFollowsForUser :many
SELECT * FROM channel_follows
WHERE user_id = $1;

-- name: GetUserFollowedChannels :many
SELECT channels.*
FROM channel_follows
JOIN channels
ON channel_follows.channel_id = channels.id
WHERE channel_follows.user_id = $1;

-- name: GetOtherChannelsForUser :many
SELECT * FROM channels
WHERE id NOT IN (
  SELECT channel_id
  FROM channel_follows
  WHERE user_id = $1
);

-- name: DeleteChannelFollow :exec
DELETE FROM channel_follows
WHERE user_id = $1 AND channel_id = $2;