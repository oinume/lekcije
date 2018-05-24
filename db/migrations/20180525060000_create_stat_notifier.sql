-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE stat_notifier (
  `datetime` datetime NOT NULL,
  `interval` tinyint(3) unsigned NOT NULL,
  `elapsed` int(10) unsigned NOT NULL,
  `user_count` int(10) unsigned NOT NULL,
  `followed_teacher_count` int(10) unsigned NOT NULL,
  PRIMARY KEY (`datetime`, `interval`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE stat_notifier;
