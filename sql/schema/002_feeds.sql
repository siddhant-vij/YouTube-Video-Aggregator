-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL UNIQUE,
  last_fetched_at TIMESTAMP
);

-- +goose Down
DROP TABLE feeds;