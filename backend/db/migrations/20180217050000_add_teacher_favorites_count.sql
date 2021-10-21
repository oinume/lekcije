-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE teacher ADD COLUMN `favorite_count` int unsigned NOT NULL DEFAULT 0 AFTER `years_of_experience`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher DROP COLUMN `favorite_count`;
