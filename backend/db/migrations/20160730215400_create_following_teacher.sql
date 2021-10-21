
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS teacher (
    id int unsigned NOT NULL,
    name varchar(255) NOT NULL DEFAULT '',
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`id`)
) ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS following_teacher (
    user_id int unsigned NOT NULL,
    teacher_id int unsigned NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`user_id`, `teacher_id`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS teacher;
DROP TABLE IF EXISTS following_teacher;
