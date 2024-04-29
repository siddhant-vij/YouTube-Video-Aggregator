-- name: InsertBookmark :exec
INSERT INTO bookmarks
  (id, created_at, updated_at, user_id, video_id)
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetAllBookmarks :many
SELECT * FROM bookmarks;

-- name: GetVideosBookmarkedByUser :many
SELECT videos.*, TRUE AS bookmark_status
FROM bookmarks
JOIN videos
ON bookmarks.video_id = videos.id
WHERE bookmarks.user_id = $1
ORDER BY DATE(published_at) DESC,
  vote_count DESC,
  view_count DESC,
  star_count DESC,
  star_rating DESC;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND video_id = $2;