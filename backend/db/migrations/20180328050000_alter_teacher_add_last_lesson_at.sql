-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE teacher ADD COLUMN `last_lesson_at` DATETIME NOT NULL AFTER `rating`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher DROP COLUMN `last_lesson_at`;
