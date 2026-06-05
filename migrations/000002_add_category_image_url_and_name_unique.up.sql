ALTER TABLE categories
ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

ALTER TABLE categories
ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES categories(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_categories_parent_id
ON categories(parent_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name_lower_unique
ON categories(LOWER(name));