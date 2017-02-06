-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
BEGIN;
UPDATE `m_country` SET name_ja='ロシア連邦' WHERE id =643;
UPDATE `m_country` SET name_ja='マケドニア共和国' WHERE id = 807;
COMMIT;

ALTER TABLE teacher
  MODIFY COLUMN `years_of_experience` TINYINT NOT NULL DEFAULT -1;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher
  MODIFY COLUMN `years_of_experience` TINYINT UNSIGNED NOT NULL DEFAULT 0;
