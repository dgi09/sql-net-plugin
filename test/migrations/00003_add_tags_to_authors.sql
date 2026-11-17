-- +goose Up
ALTER TABLE authors
ADD COLUMN tags text[];

-- +goose Down
