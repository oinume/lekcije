-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
BEGIN;

UPDATE following_teacher SET
  created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
  updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);
UPDATE lesson SET
    created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);
UPDATE teacher SET
    created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);
UPDATE user SET
    created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);
UPDATE user_api_token SET
    created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);
UPDATE user_google SET
    created_at = DATE_ADD(created_at, INTERVAL -9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL -9 HOUR);

COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

BEGIN;

UPDATE following_teacher SET
  created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
  updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);
UPDATE lesson SET
    created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);
UPDATE teacher SET
    created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);
UPDATE user SET
    created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);
UPDATE user_api_token SET
    created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);
UPDATE user_google SET
    created_at = DATE_ADD(created_at, INTERVAL 9 HOUR),
    updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR);

COMMIT;
