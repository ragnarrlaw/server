-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "citext";

-- CREATE TABLE discount (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     store_id UUID,
--     discount CITEXT,
--     applicable_tags VARCHAR(255)[], -- tags like 'holiday', 'new_year', etc.
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
    discount CITEXT,
    applicable_tags VARCHAR(255)[],
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE discount;
