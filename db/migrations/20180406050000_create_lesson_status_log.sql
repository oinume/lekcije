-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lesson
  ADD COLUMN `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT FIRST
  , DROP PRIMARY KEY
  , ADD PRIMARY KEY (id)
  , ADD UNIQUE KEY `teacher_id-datetime` (`teacher_id`, `datetime`)
;

CREATE TABLE lesson_status_log (
    `lesson_id` BIGINT UNSIGNED NOT NULL,
    `status` enum('finished','reserved','available','cancelled') COLLATE utf8mb4_bin NOT NULL,
    `created_at` DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE lesson_status_log;
ALTER TABLE `lesson`
  DROP PRIMARY KEY
  , ADD PRIMARY KEY `teacher_id-datetime` (`teacher_id`, `datetime`)
  , DROP COLUMN `id`
  , DROP KEY `teacher_id-datetime`
;
