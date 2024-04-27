-- name: InsertVideo :exec
INSERT INTO videos
  (id, created_at, updated_at, title, description, image_url, authors, published_at, url, view_count, star_rating, star_count, channel_id)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: GetAllVideos :many
SELECT * FROM videos;

-- name: GetUserVideos :many
SELECT videos.* FROM videos
JOIN channel_follows
ON videos.channel_id = channel_follows.channel_id
WHERE channel_follows.user_id = $1
ORDER BY videos.published_at DESC
LIMIT $2;

-- name: GetStatsForURL :one
SELECT view_count, star_rating, star_count
FROM videos
WHERE url = $1;

-- name: UpdateStatsForURL :exec
UPDATE videos
SET
  view_count = $2,
  star_rating = $3,
  star_count = $4
WHERE url = $1;

-- name: DeleteOldVideos :exec
DELETE FROM videos
WHERE published_at < (CURRENT_DATE - interval '1 month');
