-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

-- CREATE TABLE discount (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     store_id UUID,
--     product_id UUID,
--     discount_type VARCHAR(50) NOT NULL, -- 'percentage', 'flat', 'card_offer'
--     value NUMERIC NOT NULL, -- 10 for 10% or a flat amount
--     applicable_tags VARCHAR(255)[], -- tags like 'holiday', 'new_year', etc.
--     applicable_cards VARCHAR(255)[], -- cards eligible for the discount
--     start_date TIMESTAMP NOT NULL,
--     end_date TIMESTAMP NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY store_id REFERENCES store(id) ON DELETE CASCADE,
--     FOREIGN KEY product_id REFERENCES product(id) ON DELETE CASCADE
-- );

CREATE TABLE discount (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_id UUID NOT NULL REFERENCES store(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    discount_type VARCHAR(50) NOT NULL,
    value NUMERIC NOT NULL,
    applicable_tags VARCHAR(255)[],
    applicable_cards VARCHAR(255)[],
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE discount;
