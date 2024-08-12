-- name: CreateUser :exec
INSERT INTO users(username, password, email)
VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
