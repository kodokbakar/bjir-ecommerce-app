DROP TRIGGER IF EXISTS set_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS set_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS set_carts_updated_at ON carts;
DROP TRIGGER IF EXISTS set_products_updated_at ON products;
DROP TRIGGER IF EXISTS set_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS set_users_updated_at ON users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;