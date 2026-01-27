-- +goose Up
CREATE TABLE jobs(
    ID uuid PRIMARY KEY,
    JobID TEXT NOT NULL,
    URL TEXT UNIQUE NOT NULL,
    Fetched BOOLEAN NOT NULL DEFAULT false,
    Crawled BOOLEAN NOT NULL DEFAULT false,
    Tokenized BOOLEAN NOT NULL DEFAULT false,
    Indexed BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE jobs;