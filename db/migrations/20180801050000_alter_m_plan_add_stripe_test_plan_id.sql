-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_plan
  ADD COLUMN `stripe_test_plan_id` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL AFTER `stripe_test_product_id`;

BEGIN;
UPDATE m_plan SET stripe_test_plan_id = 'plan_Cr4O3gtBZ1Kl2m' WHERE id = 1 /* Free */;
UPDATE m_plan SET stripe_test_plan_id = 'plan_Cr4PlJaSvaf3oX' WHERE id = 2 /* Bronze */;
UPDATE m_plan SET stripe_test_plan_id = 'plan_Cr4Q7EDCRLfkjo' WHERE id = 3 /* Silver */;
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE m_plan DROP COLUMN `stripe_test_plan_id`;
