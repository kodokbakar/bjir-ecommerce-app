DROP TRIGGER IF EXISTS trg_carts_updated_at ON carts;
DROP INDEX IF EXISTS idx_carts_product_id;
DROP INDEX IF EXISTS idx_carts_user_id;
DROP TABLE IF EXISTS carts;