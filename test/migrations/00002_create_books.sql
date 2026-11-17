-- +goose Up
CREATE TABLE books (
    id           BIGSERIAL PRIMARY KEY,
    author_id    BIGINT      NOT NULL REFERENCES authors (id) ON DELETE CASCADE,
    title        TEXT        NOT NULL,
    isbn         TEXT UNIQUE,
    published_at DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX books_author_id_idx ON books (author_id);

-- +goose Down
DROP TABLE books;
