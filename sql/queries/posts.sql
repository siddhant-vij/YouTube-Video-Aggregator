-- name: InsertPost :exec
INSERT INTO posts
  (id, created_at, updated_at, title, image_title, image_url, authors, published_at, description, url, feed_id)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetAllPosts :many
SELECT * FROM posts;