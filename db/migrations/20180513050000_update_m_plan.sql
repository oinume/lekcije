-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
BEGIN;
UPDATE m_plan SET name='ブロンズ', internal_name='Bronze', notification_interval=5, show_ad=0, updated_at=NOW() WHERE id = 2;
UPDATE m_plan SET name='シルバー', internal_name='Silver', notification_interval=1, show_ad=0, updated_at=NOW() WHERE id = 3;
UPDATE m_plan SET name='ゴールド', internal_name='Gold', notification_interval=1, show_ad=0, updated_at=NOW() WHERE id = 4;
UPDATE m_plan SET name='シルバー', internal_name='Silver Beta', notification_interval=1, show_ad=0, updated_at=NOW() WHERE id = 5;
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
