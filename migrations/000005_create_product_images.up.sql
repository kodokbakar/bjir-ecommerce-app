CREATE TABLE IF NOT EXISTS product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(2048) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_images_product_id_sort_order
ON product_images(product_id, sort_order, created_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_images_one_primary_per_product
ON product_images(product_id)
WHERE is_primary = TRUE;

INSERT INTO product_images (
    product_id,
    image_url,
    sort_order,
    is_primary,
    created_at,
    updated_at
)
SELECT
    id,
    image_url,
    0,
    TRUE,
    created_at,
    updated_at
FROM products
WHERE image_url IS NOT NULL
AND image_url <> ''
AND NOT EXISTS (
    SELECT 1
    FROM product_images
    WHERE product_images.product_id = products.id
);