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
TRUNCATE TABLE users CASCADE;

-- name: GetUserPassword :one
SELECT * from users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, updated_at = $4, password = $3 WHERE id = $1 RETURNING *;