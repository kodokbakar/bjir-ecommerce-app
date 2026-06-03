CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'customer',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_role_check CHECK (role IN ('customer', 'admin'))
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(180) NOT NULL UNIQUE,
    description TEXT,
    price NUMERIC(12, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT products_price_check CHECK (price >= 0),
    CONSTRAINT products_stock_check CHECK (stock >= 0)
);

-- Catatan:
-- Tabel carts di sini berperan sebagai cart item.
-- 1 row = 1 produk dalam keranjang user.
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT carts_quantity_check CHECK (quantity > 0),
    CONSTRAINT carts_user_product_unique UNIQUE (user_id, product_id)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    shipping_address TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_status_check CHECK (
        status IN (
            'pending',
            'paid',
            'processing',
            'shipped',
            'delivered',
            'cancelled'
        )
    ),
    CONSTRAINT orders_total_amount_check CHECK (total_amount >= 0)
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_name VARCHAR(150) NOT NULL,
    quantity INTEGER NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    subtotal NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT order_items_quantity_check CHECK (quantity > 0),
    CONSTRAINT order_items_price_check CHECK (price >= 0),
    CONSTRAINT order_items_subtotal_check CHECK (subtotal >= 0)
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(120),
    amount NUMERIC(12, 2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT payments_amount_check CHECK (amount >= 0),
    CONSTRAINT payments_status_check CHECK (
        status IN (
            'pending',
            'paid',
            'failed',
            'expired',
            'refunded'
        )
    )
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_name ON users(name);

CREATE INDEX idx_categories_slug ON categories(slug);

CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_name ON products(name);

CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_product_id ON carts(product_id);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_carts_updated_at
BEFORE UPDATE ON carts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();