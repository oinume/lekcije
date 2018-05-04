-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
BEGIN;
UPDATE `m_country` SET name_ja='台湾（台湾省/中華民国）' WHERE id = 158;
UPDATE `m_country` SET name_ja='中華人民共和国' WHERE id = 156;
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
BEGIN;
UPDATE `m_country` SET name_ja='台湾（台湾省/中華民国）' WHERE id = 158;
UPDATE `m_country` SET name_ja='中華人民共和国|中国' WHERE id = 156;
COMMIT;
