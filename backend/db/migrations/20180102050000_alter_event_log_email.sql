-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE event_log_email MODIFY COLUMN
  event ENUM('click', 'delivered', 'open', 'deferred', 'dropped', 'bounce', 'block') NOT NULL;

UPDATE event_log_email SET event = 'dropped' WHERE event = '';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE event_log_email MODIFY COLUMN
  event ENUM('click', 'delivered', 'open', 'deferred', 'drop', 'bounce', 'block') NOT NULL;
