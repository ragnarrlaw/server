-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(255) NOT NULL UNIQUE,
    parent_category_id UUID DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_category_id) REFERENCES category(id) ON DELETE CASCADE
);

INSERT INTO category (category, parent_category_id) VALUES
    ('Food & Beverages', NULL),
    ('Dairy', (SELECT id FROM category WHERE category = 'Food & Beverages')),
    ('Milk', (SELECT id FROM category WHERE category = 'Dairy')),
    ('Electronics', NULL),
    ('Phones', (SELECT id FROM category WHERE category = 'Electronics'));

CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    brand VARCHAR(255),
    category_id UUID REFERENCES category(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE store_product (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     store_id UUID NOT NULL,
--     product_id UUID NOT NULL,
--     price NUMERIC NOT NULL,
--     currency VARCHAR(3) DEFAULT 'LKR',
--     unit_of_measurement VARCHAR(50) NOT NULL, -- 'kg', 'liter', 'unit'
--     stock_quantity NUMERIC NOT NULL, -- current stock quantity
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (store_id) REFERENCES store(id) ON DELETE CASCADE,
--     FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
--     UNIQUE(store_id, product_id) -- a store cannot list the same product twice
-- );

CREATE TABLE store_product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_id UUID NOT NULL,
    product_id UUID NOT NULL,
    price NUMERIC NOT NULL,
    currency VARCHAR(3) DEFAULT 'LKR',
    unit_of_measurement VARCHAR(50) NOT NULL,
    stock_quantity NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES store(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
    UNIQUE(store_id, product_id)
);

-- +goose Down
DROP TABLE store_product;
DROP TABLE product;
DROP TABLE category;
