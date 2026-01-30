-- +goose Up
ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS last_crawled_at timestamptz;

-- +goose Down
ALTER TABLE jobs
    DROP COLUMN IF EXISTS last_crawled_at;
