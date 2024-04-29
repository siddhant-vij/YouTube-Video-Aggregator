-- name: InsertVideo :exec
INSERT INTO videos
  (id, created_at, updated_at, title, description, image_url, authors, published_at, url, view_count, star_rating, star_count, vote_count, channel_id)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: GetAllVideos :many
SELECT * FROM videos;

-- name: GetUserVideos :many
SELECT v.*,
  EXISTS(
    SELECT 1 FROM bookmarks
    WHERE bookmarks.video_id = v.id
      AND bookmarks.user_id = $1
  ) AS bookmark_status
FROM (
  SELECT * FROM videos
  WHERE channel_id IN (
    SELECT channel_id FROM channel_follows
    WHERE user_id = $1
  )
  ORDER BY DATE(published_at) DESC,
    vote_count DESC,
    view_count DESC,
    star_count DESC,
    star_rating DESC
  LIMIT $2
) v;

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

-- name: UpvoteVideo :exec
UPDATE videos
SET
  vote_count = vote_count + 1
WHERE id = $1;

-- name: DownvoteVideo :exec
UPDATE videos
SET
  vote_count = vote_count - 1
WHERE id = $1;
