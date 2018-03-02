-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE teacher
  ADD COLUMN `review_count` int unsigned NOT NULL DEFAULT 0 AFTER `favorite_count`,
  ADD COLUMN `rating` DECIMAL(2, 1) DEFAULT 0.0 AFTER `review_count`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher DROP COLUMN `review_count`, DROP COLUMN `rating`;
