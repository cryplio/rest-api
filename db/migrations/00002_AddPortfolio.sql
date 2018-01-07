-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE user_portfolios (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  name VARCHAR NOT NULL UNIQUE CHECK (length(email) < 50),
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
