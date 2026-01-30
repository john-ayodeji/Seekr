-- +goose Up
ALTER TABLE parsed_html
    ADD COLUMN search_vector tsvector;

CREATE INDEX parsed_html_search_idx
    ON parsed_html
        USING GIN (search_vector);

ALTER TABLE parsed_html
    ADD COLUMN indexed_at timestamptz DEFAULT now();