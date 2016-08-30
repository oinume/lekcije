-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS lesson (
  `teacher_id` int unsigned NOT NULL,
  `datetime` datetime NOT NULL,
  `status` enum('finished', 'reserved', 'available', 'cancelled') NOT NULL,
  PRIMARY KEY (`teacher_id`,`datetime`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS lesson;
