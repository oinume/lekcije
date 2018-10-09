-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'チェコ共和国' WHERE id = 203;
UPDATE m_country SET name_ja = 'ジョージア' WHERE id = 268;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
