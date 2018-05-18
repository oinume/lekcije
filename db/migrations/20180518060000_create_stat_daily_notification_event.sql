-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE stat_daily_notification_event (
 `date` date NOT NULL,
  `event` enum('click','delivered','open','deferred','dropped','bounce','block') COLLATE utf8mb4_bin NOT NULL,
  `count` int(10) unsigned NOT NULL,
  `uu_count` int(10) unsigned NOT NULL,
  PRIMARY KEY (`date`,`event`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

INSERT INTO stat_daily_notification_event SELECT * FROM stat_new_lesson_notifier;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE stat_daily_notification_event;
