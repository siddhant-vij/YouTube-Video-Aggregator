-- +goose Up
CREATE TABLE videos (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  image_url VARCHAR(255) NOT NULL,
  authors VARCHAR(255) NOT NULL,
  published_at TIMESTAMP NOT NULL,
  url VARCHAR(255) NOT NULL UNIQUE,
  view_count VARCHAR(10) NOT NULL,
  star_rating DECIMAL NOT NULL,
  star_count VARCHAR(10) NOT NULL,
  channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE videos;