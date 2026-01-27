-- +goose Up
CREATE TABLE parsed_html(
    id uuid PRIMARY KEY ,
    jobId TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    headings TEXT NOT NULL,
    paragraphs TEXT NOT NULL,
    links TEXT NOT NULL
);

-- +goose Down
DROP TABLE parsed_html;