-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS m_country (
  `id` SMALLINT unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `name_ja` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

ALTER TABLE teacher
  ADD COLUMN `country_id` SMALLINT NOT NULL DEFAULT 0 AFTER `name`,
  ADD COLUMN `gender` ENUM('female', 'male', 'other') NOT NULL DEFAULT 'other' AFTER `country_id`,
  ADD COLUMN `birthday` DATETIME NOT NULL AFTER `gender`,
  ADD COLUMN `years_of_experience` TINYINT UNSIGNED NOT NULL DEFAULT 0 AFTER `birthday`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS m_country;
ALTER TABLE teacher
  DROP COLUMN `country_id`,
  DROP COLUMN `gender`,
  DROP COLUMN `birthday`,
  DROP COLUMN ``years_of_experience;
