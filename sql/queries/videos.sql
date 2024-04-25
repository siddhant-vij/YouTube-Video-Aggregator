-- name: InsertVideo :exec
INSERT INTO videos
  (id, created_at, updated_at, title, description, image_url, authors, published_at, url, view_count, star_rating, star_count, channel_id)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: GetAllVideos :many
SELECT * FROM videos;