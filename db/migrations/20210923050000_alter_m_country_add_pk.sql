-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_country ADD PRIMARY KEY (`id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
