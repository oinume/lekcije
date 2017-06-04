-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO plan VALUES (5, 'Plus', 'Plus Beta', 0, 1, 0, '2017-06-04 00:00:00', '2017-06-04 00:00:00');
CREATE TABLE IF NOT EXISTS m_plan LIKE plan;
INSERT INTO m_plan SELECT * FROM plan;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS m_plan;
DELETE FROM plan WHERE id = 5;
