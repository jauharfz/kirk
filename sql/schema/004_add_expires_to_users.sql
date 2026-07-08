-- +goose up
ALTER TABLE users ADD COLUMN expires_in_seconds INTEGER NOT NULL DEFAULT 3600;

-- +goose down
ALTER TABLE users DROP COLUMN expires_in_seconds;
