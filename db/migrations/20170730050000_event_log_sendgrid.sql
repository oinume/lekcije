-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS event_log_email (
  datetime DATETIME NOT NULL,
  event ENUM('click', 'delivered', 'open', 'deferred', 'drop', 'bounce', 'block') NOT NULL,
  email_type ENUM('new_lesson') NOT NULL,
  user_id int(10) unsigned NOT NULL,
  user_agent VARCHAR(255) DEFAULT NULL,
  teacher_ids TEXT DEFAULT NULL,
  url VARCHAR(255) DEFAULT NULL,
  KEY (`datetime`, `event`),
  KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS event_log_email;
