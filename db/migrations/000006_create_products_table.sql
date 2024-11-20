-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    brand VARCHAR(255),
    category_id UUID REFERENCES category(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE store_product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_id UUID NOT NULL REFERENCES store(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    price NUMERIC NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'LKR',
    unit_of_measurement VARCHAR(50) NOT NULL, -- 'kg', 'liter', 'unit'
    stock_quantity NUMERIC NOT NULL, -- current stock quantity
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(store_id, product_id) -- a store cannot list the same product twice
);

CREATE TABLE discount (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_id UUID REFERENCES store(id) ON DELETE CASCADE,
    product_id UUID REFERENCES product(id) ON DELETE CASCADE,
    discount_type VARCHAR(50) NOT NULL, -- 'percentage', 'flat', 'card_offer'
    value NUMERIC NOT NULL, -- 10 for 10% or a flat amount
    applicable_tags VARCHAR(255)[], -- tags like 'holiday', 'new_year', etc.
    applicable_cards VARCHAR(255)[], -- cards eligible for the discount
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(255) NOT NULL UNIQUE,
    parent_category_id UUID REFERENCES category(id) DEFAULT NULL ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO category (id, category, parent_category_id) VALUES
    ('Food & Beverages', NULL),
    ('Dairy', (SELECT id FROM product_categories WHERE category = 'Food & Beverages')),
    ('Milk', (SELECT id FROM product_categories WHERE category = 'Dairy')),
    ('Electronics', NULL),
    ('Phones', (SELECT id FROM product_categories WHERE category = 'Electronics'));


-- +goose Down
DROP TABLE category;
DROP TABLE product_category;
DROP TABLE discount;
DROP TABLE store_products;
DROP TABLE stores;
