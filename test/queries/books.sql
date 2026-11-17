-- name: GetBook :one
SELECT b.id, a.id as author_id, a.name as author_name,a.bio as author_bio FROM books b
LEFT JOIN authors a on b.author_id = a.id
WHERE b.id = $1 
FOR UPDATE OF b
LIMIT 1;

-- name: ListBooksByAuthor :many
SELECT * FROM books
WHERE author_id = $1
ORDER BY published_at DESC NULLS LAST;

-- name: CreateBook :one
INSERT INTO books (author_id, title, isbn, published_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: ListBooksWithAuthor :many
SELECT b.id, b.title, b.isbn, b.published_at, a.id AS author_id, a.name AS author_name
FROM books b
JOIN authors a ON a.id = b.author_id
ORDER BY b.title;

-- name: CountBooksPerAuthor :many
SELECT a.id, a.name, COUNT(b.id) AS count
FROM authors a
LEFT JOIN books b ON b.author_id = a.id
GROUP BY a.id, a.name
ORDER BY count DESC;
