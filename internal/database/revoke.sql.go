// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: revoke.sql

package database

import (
	"context"
)

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP WHERE token = $1 RETURNING token, created_at, updated_at, user_id, expired_at, revoked_at
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiredAt,
		&i.RevokedAt,
	)
	return i, err
}
