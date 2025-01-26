-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP WHERE token = $1 RETURNING *;