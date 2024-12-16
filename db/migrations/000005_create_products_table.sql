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

CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(255) DEFAULT '',
    name VARCHAR(1000) NOT NULL,
    description TEXT DEFAULT '',
    brand VARCHAR(1000) DEFAULT '',
    brand_tags TEXT DEFAULT '',
    category_id UUID DEFAULT NULL,
    labels TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    product_quantity VARCHAR(255) DEFAULT '',
    serving_size VARCHAR(255) DEFAULT '',
    unit_of_measure VARCHAR(20) NOT NULL CHECK (unit_of_measure IN ('pcs', 'g', 'kg', 'ml', 'l')),
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
    price_per_unit NUMERIC NOT NULL,
    currency VARCHAR(3) DEFAULT 'LKR',
    listed_unit_of_measure VARCHAR(20) NOT NULL CHECK (listed_unit_of_measure IN ('pcs', 'g', 'kg', 'ml', 'l')),
    stock_quantity VARCHAR(20) NOT NULL CHECK (stock_quantity IN ('LIMITED_QUANTITY', 'AVAILABLE', 'NOT_AVAILABLE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES store(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
    UNIQUE(store_id, product_id)
);

CREATE TYPE recommendation_item_t AS (
    product_id UUID,
    price_per_unit NUMERIC,
    listed_uom TEXT,
    stock_quantity INT,
    product_name TEXT,
    product_image_url TEXT,
    standard_uom TEXT
);

-- +goose Down
DROP TABLE store_product;
DROP TABLE product;
DROP TYPE IF EXISTS recommendation_item_t;
DROP TABLE category;
