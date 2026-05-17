-- +goose Up
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    account_number VARCHAR(30) NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'NPR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);

-- +goose Down
DROP TABLE IF EXISTS accounts;