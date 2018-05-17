-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE m_plan
  ADD COLUMN `stripe_test_product_id` VARCHAR(255) DEFAULT NULL AFTER `internal_name`;

BEGIN;
UPDATE m_plan SET stripe_test_product_id='prod_Cr4OmqTNz986CK' WHERE id = 1;
UPDATE m_plan SET stripe_test_product_id='prod_Cr4Ofvm4M2uxDd' WHERE id = 2;
UPDATE m_plan SET stripe_test_product_id='prod_Cr4PyhkwibtGdF' WHERE id = 3;
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE m_plan DROP COLUMN `stripe_test_product_id`;
