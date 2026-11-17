-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = @id
LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
WHERE 
	(sqlc.narg(min_id)::bigint IS NULL OR id > sqlc.narg(min_id)) AND
    (sqlc.narg(min_created_at)::timestamptz IS NULL OR created_at > sqlc.narg(min_created_at))
ORDER BY
  CASE WHEN sqlc.arg(sort_by)::text = 'id' THEN id END ASC,
  CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
  CASE WHEN sqlc.arg(sort_by)::text = 'created_at' THEN created_at END DESC;

-- name: ListAuthorsWithTags :many
SELECT * FROM authors
WHERE tags IS NOT NULL AND tags @> $1
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio , tags, created_at)
VALUES ($1, $2, $3, now())
RETURNING id;

-- name: UpdateAuthor :exec
UPDATE authors
SET name = $2, bio = $3
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

