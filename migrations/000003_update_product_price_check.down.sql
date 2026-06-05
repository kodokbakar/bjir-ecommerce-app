ALTER TABLE products
DROP CONSTRAINT IF EXISTS products_price_positive_check;

ALTER TABLE products
ADD CONSTRAINT products_price_check CHECK (price >= 0);