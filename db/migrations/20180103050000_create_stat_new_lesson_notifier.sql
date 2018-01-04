-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS stat_new_lesson_notifier (
   `date` date NOT NULL,
   `event` enum('click','delivered','open','deferred','drop','bounce','block') NOT NULL,
   `count` int unsigned NOT NULL,
   `uu_count` int unsigned NOT NULL,
   PRIMARY KEY (`date`, `event`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE stat_new_lesson_notifier;
