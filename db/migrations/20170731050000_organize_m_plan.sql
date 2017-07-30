-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE IF EXISTS plan;

BEGIN;
UPDATE m_plan SET name='フリー', notification_interval=10, updated_at='2017-07-30 13:43:31' WHERE id = 1;
UPDATE m_plan SET name='プラス', internal_name='Plus', price=400, notification_interval=1, updated_at='2017-07-30 13:43:31' WHERE id = 2;
UPDATE m_plan SET name='プロ', internal_name='Pro', price=900, notification_interval=1, updated_at='2017-07-30 13:43:31' WHERE id = 3;
UPDATE m_plan SET name='プラス', internal_name='Plus Beta', price=0, updated_at='2017-07-30 13:43:31' WHERE id = 5;
COMMIT;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
CREATE TABLE IF NOT EXISTS `plan` (
  `id` tinyint(3) unsigned NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `internal_name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `price` mediumint(9) NOT NULL,
  `notification_interval` tinyint(3) unsigned NOT NULL,
  `show_ad` tinyint(1) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS m_plan;
CREATE TABLE `m_plan` (
  `id` tinyint(3) unsigned NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `internal_name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `price` mediumint(9) NOT NULL,
  `notification_interval` tinyint(3) unsigned NOT NULL,
  `show_ad` tinyint(1) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
INSERT INTO `m_plan` VALUES (1,'Free','Free',0,20,1,'2016-10-16 00:00:00','2016-10-16 00:00:00'),(2,'Standard','Standard',300,10,0,'2016-10-16 00:00:00','2016-10-16 00:00:00'),(3,'Plus','Plus',500,5,0,'2016-10-16 00:00:00','2016-10-16 00:00:00'),(4,'Free','Free Beta',0,10,1,'2016-10-16 00:00:00','2016-10-16 00:00:00'),(5,'Plus','Plus Beta',0,1,0,'2017-06-04 00:00:00','2017-06-04 00:00:00');
