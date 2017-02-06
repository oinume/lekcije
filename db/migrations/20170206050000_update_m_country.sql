-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
BEGIN;
UPDATE `m_country` SET name_ja='ロシア連邦' WHERE id =643;
UPDATE `m_country` SET name_ja='マケドニア共和国' WHERE id = 807;
COMMIT;

/*
ALTER TABLE teacher
  ADD COLUMN `country_id` SMALLINT NOT NULL DEFAULT 0 AFTER `name`,
  ADD COLUMN `gender` ENUM('female', 'male', 'other') NOT NULL DEFAULT 'other' AFTER `country_id`,
  ADD COLUMN `birthday` DATE NOT NULL AFTER `gender`,
  ADD COLUMN `years_of_experience` TINYINT UNSIGNED NOT NULL DEFAULT 0 AFTER `birthday`;
*/

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
/*
ALTER TABLE teacher
  DROP COLUMN `country_id`,
  DROP COLUMN `gender`,
  DROP COLUMN `birthday`,
  DROP COLUMN `years_of_experience`;
*/
