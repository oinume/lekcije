-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE event_log_email
    ADD COLUMN `id` int unsigned NOT NULL AUTO_INCREMENT FIRST,
    ADD PRIMARY KEY (id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
