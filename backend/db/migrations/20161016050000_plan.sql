-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS plan (
  `id` tinyint unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `internal_name` varchar(255) NOT NULL,
  `price` mediumint NOT NULL,
  `notification_interval` tinyint unsigned NOT NULL,
  `show_ad` tinyint(1) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

INSERT INTO plan VALUES (1, 'Free', 'Free', 0, 20, 1, '2016-10-16 00:00:00', '2016-10-16 00:00:00');
INSERT INTO plan VALUES (2, 'Standard', 'Standard', 300, 10, 0, '2016-10-16 00:00:00', '2016-10-16 00:00:00');
INSERT INTO plan VALUES (3, 'Plus', 'Plus', 500, 5, 0, '2016-10-16 00:00:00', '2016-10-16 00:00:00');

ALTER TABLE user ADD COLUMN `plan_id` tinyint unsigned NOT NULL DEFAULT 1 AFTER `email_verified`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS plan;
ALTER TABLE user DROP COLUMN `plan_id`;
