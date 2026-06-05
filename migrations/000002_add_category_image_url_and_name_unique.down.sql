DROP INDEX IF EXISTS idx_categories_name_lower_unique;
DROP INDEX IF EXISTS idx_categories_parent_id;

ALTER TABLE categories
DROP COLUMN IF EXISTS parent_id;

ALTER TABLE categories
DROP COLUMN IF EXISTS image_url;