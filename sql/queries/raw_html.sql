-- name: AddHtml :one
INSERT INTO raw_html(id, jobid, url, html)
VALUES ($1, $2, $3, $4)
RETURNING *;