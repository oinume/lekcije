-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS event_log_email (
  datetime DATETIME NOT NULL,
  event ENUM('click', 'delivered', 'open', 'deferred', 'drop', 'bounce', 'block') NOT NULL,
  email_type ENUM('new_lesson_notifier', 'follow_reminder') NOT NULL,
  user_id int(10) unsigned NOT NULL,
  user_agent VARCHAR(255) NOT NULL DEFAULT '',
  teacher_ids TEXT NOT NULL,
  url VARCHAR(255) NOT NULL DEFAULT '',
  KEY (`datetime`, `event`),
  KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=COMPRESSED;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS event_log_email;
