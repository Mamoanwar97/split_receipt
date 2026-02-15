-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE friends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_name TEXT NOT NULL,
    date_time TIMESTAMPTZ NOT NULL,
    total NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    total_price NUMERIC(10,2) NOT NULL,
    price_per_unit NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE meal_friends (
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    friend_id UUID NOT NULL REFERENCES friends(id) ON DELETE CASCADE,
    PRIMARY KEY (meal_id, friend_id)
);

CREATE TABLE receipt_fixed_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL
);

CREATE TABLE receipt_percentage_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    percentage NUMERIC(5,2) NOT NULL,
    cap_amount NUMERIC(10,2)
);

-- +goose Down
DROP TABLE IF EXISTS receipt_percentage_fees;
DROP TABLE IF EXISTS receipt_fixed_fees;
DROP TABLE IF EXISTS meal_friends;
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS receipts;
DROP TABLE IF EXISTS friends;
