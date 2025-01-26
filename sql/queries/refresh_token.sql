-- name: GetUserFromRefreshToken :one
SELECT
    users.*
FROM
    users
    JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE
    refresh_tokens.token = $1
    AND refresh_tokens.expired_at > CURRENT_TIMESTAMP;


-- name: CreateRefreshToken :one
INSERT INTO
    refresh_tokens(
        token,
        created_at,
        updated_at,
        user_id,
        expired_at
    )
    VALUES
    ($1, $2, $3, $4, $5) RETURNING *;