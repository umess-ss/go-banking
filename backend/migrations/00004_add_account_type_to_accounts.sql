-- +goose Up
ALTER TABLE accounts
ADD COLUMN account_type VARCHAR(20) NOT NULL DEFAULT 'savings';

-- +goose Down
ALTER TABLE accounts
DROP COLUMN account_type;
