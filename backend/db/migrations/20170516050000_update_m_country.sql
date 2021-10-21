-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE m_country SET name_ja='南アフリカ' WHERE id = 710;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
UPDATE m_country SET name_ja='南アフリカ共和国|南アフリカ' WHERE id = 710;
