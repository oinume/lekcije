-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_plan
  ADD COLUMN `max_following_teacher` TINYINT unsigned NOT NULL DEFAULT 0 AFTER `show_ad`;

BEGIN;
UPDATE m_plan SET max_following_teacher = 10 WHERE id = 1 /* Free */;
UPDATE m_plan SET max_following_teacher = 15 WHERE id = 2 /* Bronze */;
UPDATE m_plan SET max_following_teacher = 15 WHERE id = 3 /* Silver */;
UPDATE m_plan SET max_following_teacher = 20 WHERE id = 4 /* Gold */;
UPDATE m_plan SET max_following_teacher = 15 WHERE id = 5 /* Silver Beta */;
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE m_plan DROP COLUMN `max_following_teacher`;
