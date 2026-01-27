-- +goose Up
CREATE TABLE raw_html(
    ID uuid PRIMARY KEY,
    JobID TEXT NOT NULL,
    URL TEXT UNIQUE NOT NULL,
    HTML TEXT NOT NULL
);

-- +goose Down
DROP TABLE raw_html;