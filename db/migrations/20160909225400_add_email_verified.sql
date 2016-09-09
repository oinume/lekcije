-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE user ADD COLUMN `email_verified` tinyint(1) unsigned NOT NULL DEFAULT 0 AFTER `email`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE user DROP COLUMN `email_verified`;
