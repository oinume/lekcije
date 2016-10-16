-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO plan VALUES (4, 'Free', 'Free Beta', 0, 10, 1, '2016-10-16 00:00:00', '2016-10-16 00:00:00');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM plan WHERE id = 4;
