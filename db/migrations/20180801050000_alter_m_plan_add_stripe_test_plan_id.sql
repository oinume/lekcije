-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_plan
  ADD COLUMN `stripe_test_plan_id` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL AFTER `stripe_test_product_id`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE m_plan DROP COLUMN `stripe_test_plan_id`;
