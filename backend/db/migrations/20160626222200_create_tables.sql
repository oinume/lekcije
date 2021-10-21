
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS user (
    id int unsigned NOT NULL AUTO_INCREMENT,
    name varchar(50) NOT NULL,
    email varchar(255) NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`email`)
) ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS auth_google (
    user_id int unsigned NOT NULL DEFAULT 0,
    access_token varchar(255) NOT NULL,
    id_token varchar(1024) NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY (`access_token`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS user;
