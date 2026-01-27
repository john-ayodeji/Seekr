-- name: AddParsedHTML :one
INSERT INTO parsed_html(
    id , jobId, url, title, description, headings, paragraphs, links
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;