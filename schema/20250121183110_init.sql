-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE, -- Username must be unique
    email VARCHAR(255) NOT NULL UNIQUE,    -- Email must be unique
    provider TEXT NOT NULL CHECK (provider IN ('google', 'facebook', 'local')), -- Enum-like constraint
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    content TEXT NOT NULL,
    user_id INT NOT NULL,
    visited_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add indexes to frequently queried columns
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_links_user_id ON links(user_id);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
