-- name: AddParsedHTML :one
INSERT INTO parsed_html(
    id , jobId, url, title, description, headings, paragraphs, links
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ReindexParsedHTMLByID :exec
UPDATE parsed_html
SET search_vector =
        setweight(to_tsvector('english', title), 'A') ||
        setweight(to_tsvector('english', headings), 'A') ||
        setweight(to_tsvector('english', description), 'B') ||
        setweight(to_tsvector('english', paragraphs), 'C')
WHERE id = $1;

-- name: SearchParsedHTML :many
SELECT
    id,
    url,
    title,
    GREATEST(
            (ts_rank_cd(search_vector, query) *
             (1 + 1 / (EXTRACT(EPOCH FROM (now() - indexed_at)) / 86400 + 1))),
            similarity(title, $1),
            similarity(paragraphs, $1)
    )::double precision AS rank
FROM parsed_html,
     to_tsquery('english', $2) query
WHERE search_vector @@ query
   OR title % $1
   OR paragraphs % $1
ORDER BY rank DESC, id
LIMIT $3
    OFFSET $4;



-- name: SearchParsedHTMLWithSnippet :many
SELECT
    id,
    url,
    title,
    ts_headline(
            'english',
            paragraphs,
            query,
            'MaxWords=30, MinWords=15'
    ) AS snippet,
    GREATEST(
            (ts_rank_cd(search_vector, query) *
             (1 + 1 / (EXTRACT(EPOCH FROM (now() - indexed_at)) / 86400 + 1))),
            similarity(title, $1),
            similarity(paragraphs, $1)
    )::double precision AS rank
FROM parsed_html,
     to_tsquery('english', $2) query
WHERE search_vector @@ query
   OR title % $1
   OR paragraphs % $1
ORDER BY rank DESC, id
LIMIT $3
    OFFSET $4;

