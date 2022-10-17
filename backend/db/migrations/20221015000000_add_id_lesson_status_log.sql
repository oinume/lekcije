-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lesson_status_log
    ADD COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT FIRST,
    ADD PRIMARY KEY (id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
