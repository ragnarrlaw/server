-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE store (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_username VARCHAR(255) NOT NULL UNIQUE,
    store_name VARCHAR(255) NOT NULL,
    store_address VARCHAR(255),
    store_email VARCHAR(255),
    store_contact_number VARCHAR(20),
    store_location_point GEOGRAPHY(Point, 4326),
    store_web_url VARCHAR(255),
    password_digest VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE store;
