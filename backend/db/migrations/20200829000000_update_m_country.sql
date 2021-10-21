-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja = 'タイ' where id = 764;
UPDATE m_country SET name_ja = 'ベネズエラ・ボリバル共和国' WHERE id = 862;
UPDATE m_country SET name_ja = 'モンゴル' WHERE id = 496;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
