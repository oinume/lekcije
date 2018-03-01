-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE teacher ADD COLUMN `rating` DECIMAL(2, 1) DEFAULT 0.0 AFTER `favorite_count`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher DROP COLUMN `rating`;
