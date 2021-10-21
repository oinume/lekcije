
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS user_api_token (
    token varchar(128) NOT NULL,
    user_id int unsigned NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`token`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS user_api_token;
