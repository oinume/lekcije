-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'セントビンセントおよびグレナディーン諸島' WHERE id = 670;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
