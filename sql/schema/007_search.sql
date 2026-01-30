-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX parsed_html_title_trgm_idx
    ON parsed_html USING GIN (title gin_trgm_ops);

CREATE INDEX parsed_html_paragraphs_trgm_idx
    ON parsed_html USING GIN (paragraphs gin_trgm_ops);