-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE user DROP COLUMN `raw_email`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE user ADD COLUMN `raw_email` varchar(255) CHARACTER SET utf8mb4 NOT NULL;
