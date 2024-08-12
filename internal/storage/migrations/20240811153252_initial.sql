-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY ,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password bytea NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users
