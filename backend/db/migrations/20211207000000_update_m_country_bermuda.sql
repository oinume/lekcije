-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'バミューダ', updated_at = NOW() where id = 60;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
