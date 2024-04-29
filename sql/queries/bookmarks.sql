-- name: InsertBookmark :exec
INSERT INTO bookmarks
  (id, created_at, updated_at, user_id, video_id)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllBookmarks :many
SELECT * FROM bookmarks;

-- name: GetVideoIdsBookmarkedByUser :many
SELECT video_id FROM bookmarks
WHERE user_id = $1;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND video_id = $2;