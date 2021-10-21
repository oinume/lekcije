-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS stat_daily_user_notification_event (
    `date` date NOT NULL,
    `user_id` int unsigned NOT NULL,
    `event` enum('open') NOT NULL,
    `count` int unsigned NOT NULL,
    PRIMARY KEY (`date`, `user_id`, `event`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE stat_daily_user_notification_event;
