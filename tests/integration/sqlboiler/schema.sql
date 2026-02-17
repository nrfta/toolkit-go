-- Products: Main test table with all comparator types
CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,  -- enum: active, inactive, deleted
    price DECIMAL(10,2),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    tags TEXT[],  -- Array column
    categories TEXT[]  -- Array column
);

-- Users: Test nullable fields and authorization
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT,  -- nullable
    role TEXT NOT NULL,
    organization_id TEXT,  -- nullable
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Orders: Test date comparators and relationships
CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    product_id TEXT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    order_date TIMESTAMP NOT NULL DEFAULT NOW(),
    shipped_date TIMESTAMP,  -- nullable
    status TEXT NOT NULL
);

-- Tags: Additional array testing
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    product_ids TEXT[],  -- Array of product IDs
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for query performance
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_order_date ON orders(order_date);
