-- +goose Up
CREATE TABLE channels (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL UNIQUE,
  last_fetched_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE channels;