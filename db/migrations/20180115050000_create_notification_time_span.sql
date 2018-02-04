-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS notification_time_span (
    `user_id` int unsigned NOT NULL,
    `number` tinyint unsigned NOT NULL,
    `from_time` time NOT NULL,
    `to_time` time NOT NULL,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`user_id`, `number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE notification_time_span;
