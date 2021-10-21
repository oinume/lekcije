-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE teacher MODIFY COLUMN `rating` decimal(3,2) DEFAULT '0.00';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE teacher MODIFY COLUMN `rating` decimal(2,1) DEFAULT '0.0';
