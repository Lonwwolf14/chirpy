-- +goose Up
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        created_at TIMESTAMP,
        updated_at TIMESTAMP
    );

ALTER TABLE users
ADD COLUMN password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
DROP TABLE users;