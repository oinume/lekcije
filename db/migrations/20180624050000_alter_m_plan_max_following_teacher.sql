-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_plan
  ADD COLUMN `max_following_teacher` TINYINT unsigned NOT NULL DEFAULT 0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE event_log_email
  DROP COLUMN `max_following_teacher`;
