-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    created_at,
    updated_at,
    password
)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteUsers :many
TRUNCATE TABLE users;