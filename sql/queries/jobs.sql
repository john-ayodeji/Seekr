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

-- name: MarkTokenizedAndIndexed :exec
UPDATE jobs
SET tokenized = true, indexed = true
WHERE url = $1;