-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE user ADD COLUMN `followed_teacher_at` datetime DEFAULT NULL AFTER `plan_id`;
UPDATE user SET followed_teacher_at = created_at WHERE followed_teacher_at IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE user DROP COLUMN `followed_teacher_at`;
