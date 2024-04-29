// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: bookmarks.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const deleteBookmark = `-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND video_id = $2
`

type DeleteBookmarkParams struct {
	UserID  uuid.UUID
	VideoID uuid.UUID
}

func (q *Queries) DeleteBookmark(ctx context.Context, arg DeleteBookmarkParams) error {
	_, err := q.db.ExecContext(ctx, deleteBookmark, arg.UserID, arg.VideoID)
	return err
}

const getAllBookmarks = `-- name: GetAllBookmarks :many
SELECT id, created_at, updated_at, user_id, video_id FROM bookmarks
`

func (q *Queries) GetAllBookmarks(ctx context.Context) ([]Bookmark, error) {
	rows, err := q.db.QueryContext(ctx, getAllBookmarks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bookmark
	for rows.Next() {
		var i Bookmark
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.VideoID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVideosBookmarkedByUser = `-- name: GetVideosBookmarkedByUser :many
SELECT videos.id, videos.created_at, videos.updated_at, videos.title, videos.description, videos.image_url, videos.authors, videos.published_at, videos.url, videos.view_count, videos.star_rating, videos.star_count, videos.channel_id, TRUE AS bookmark_status
FROM bookmarks
JOIN videos
ON bookmarks.video_id = videos.id
WHERE bookmarks.user_id = $1
`

type GetVideosBookmarkedByUserRow struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Title          string
	Description    string
	ImageUrl       string
	Authors        string
	PublishedAt    time.Time
	Url            string
	ViewCount      string
	StarRating     string
	StarCount      string
	ChannelID      string
	BookmarkStatus bool
}

func (q *Queries) GetVideosBookmarkedByUser(ctx context.Context, userID uuid.UUID) ([]GetVideosBookmarkedByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getVideosBookmarkedByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVideosBookmarkedByUserRow
	for rows.Next() {
		var i GetVideosBookmarkedByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.ImageUrl,
			&i.Authors,
			&i.PublishedAt,
			&i.Url,
			&i.ViewCount,
			&i.StarRating,
			&i.StarCount,
			&i.ChannelID,
			&i.BookmarkStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertBookmark = `-- name: InsertBookmark :exec
INSERT INTO bookmarks
  (id, created_at, updated_at, user_id, video_id)
VALUES
  ($1, $2, $3, $4, $5)
`

type InsertBookmarkParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	VideoID   uuid.UUID
}

func (q *Queries) InsertBookmark(ctx context.Context, arg InsertBookmarkParams) error {
	_, err := q.db.ExecContext(ctx, insertBookmark,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.VideoID,
	)
	return err
}
