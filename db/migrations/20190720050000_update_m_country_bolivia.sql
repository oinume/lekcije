-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'ボリビア多民族国' where id = 68;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
