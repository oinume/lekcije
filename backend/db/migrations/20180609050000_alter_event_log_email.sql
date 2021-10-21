-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE event_log_email
  MODIFY COLUMN `email_type` enum('new_lesson_notifier','follow_reminder', 'registration') COLLATE utf8mb4_bin NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE event_log_email
  MODIFY COLUMN `email_type` enum('new_lesson_notifier','follow_reminder') COLLATE utf8mb4_bin NOT NULL;
