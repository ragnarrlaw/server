-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE store (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) DEFAULT '',
    email VARCHAR(255) DEFAULT '',
    contact_number VARCHAR(20) DEFAULT '',
    discounts VARCHAR(1000)[] DEFAULT '{}',
    location_point GEOGRAPHY(Point, 4326) DEFAULT NULL,
    web_url VARCHAR(255) DEFAULT '',
    password_digest VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE store;
