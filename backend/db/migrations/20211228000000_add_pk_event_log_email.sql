-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE event_log_email ADD PRIMARY KEY (`datetime`, `event`, `email_type`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
