-- name: CreateJob :exec
INSERT INTO jobs(
    ID, JOBID, URL
) VALUES (
    $1, $2, $3
);

-- name: MarkFetched :exec
UPDATE jobs
SET fetched = true
WHERE url = $1;

-- name: MarkCrawled :exec
UPDATE jobs
SET crawled = true
WHERE url = $1;

-- name: MarkTokenized :exec
UPDATE jobs
SET tokenized = true
WHERE url = $1;