-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'イラン' where id = 364;
UPDATE m_country SET name_ja = '北マケドニア共和国' WHERE id = 807;
UPDATE m_country SET name_ja = 'エスワティニ王国' WHERE id = 748;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
