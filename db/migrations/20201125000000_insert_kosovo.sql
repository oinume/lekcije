-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO m_country VALUES (999, 'Kosovo', 'コソボ', NOW(), NOW());

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

