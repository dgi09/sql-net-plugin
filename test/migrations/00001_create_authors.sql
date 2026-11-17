-- +goose Up
CREATE TABLE authors (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    bio        TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE authors;
