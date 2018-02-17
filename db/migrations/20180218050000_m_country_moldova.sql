-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'モルドバ共和国' WHERE id = 498;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
