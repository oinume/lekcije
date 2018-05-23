-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE user ADD COLUMN open_notification_at datetime DEFAULT NULL AFTER followed_teacher_at;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE user DROP COLUMN open_notification_at;
