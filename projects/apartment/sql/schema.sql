-- Enable UUID generator extension (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Drop tables if they already exist
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS users_apartments;
DROP TABLE IF EXISTS apartments;
DROP TABLE IF EXISTS users;

-- Drop ENUM type if it exists
DROP TYPE IF EXISTS bill_type;

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE INDEX idx_users_email ON users(email);

-- Create apartments table
CREATE TABLE apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    name TEXT,
    address TEXT,
    unit_number INTEGER,
    admin_id UUID NOT NULL, 
    FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create junction table for many-to-many relationship between users and apartments
CREATE TABLE users_apartments (
    user_id UUID NOT NULL,
    apartment_id UUID NOT NULL,
    PRIMARY KEY (user_id, apartment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE
);


-- Create enum type for bills
CREATE TYPE bill_type AS ENUM ('electricity', 'water', 'gas');

-- Create bills table
CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    name TEXT,
    bill_type bill_type NOT NULL,
    bill_id INTEGER NOT NULL, 
    amount INTEGER,
    due_date DATE NOT NULL,
    image_id UUID,
    apartment_id UUID NOT NULL,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE
);

