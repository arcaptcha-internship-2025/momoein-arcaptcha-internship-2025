-- Enable UUID generator extension (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Drop tables if they already exist
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS users_apartments;
DROP TABLE IF EXISTS apartments;
DROP TABLE IF EXISTS users;

-- Drop ENUM type if it exists
DROP TYPE IF EXISTS bill_type;
DROP TYPE IF EXISTS invite_status_type;
DROP TYPE IF EXISTS bill_status_type;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create apartments table
CREATE TABLE IF NOT EXISTS apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    unit_number INTEGER NOT NULL,
    admin_id UUID NOT NULL, 
    FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create enum type for invite status
CREATE TYPE IF NOT EXISTS invite_status_type AS ENUM ('pending', 'accepted', 'declined');

-- Create junction table for many-to-many relationship between users and apartments
CREATE TABLE IF NOT EXISTS users_apartments (
    user_id UUID NOT NULL,
    apartment_id UUID NOT NULL,
    PRIMARY KEY (user_id, apartment_id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    invite_status invite_status_type NOT NULL DEFAULT 'pending',
    invite_token TEXT UNIQUE NOT NULL,
    invite_expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- Create enum type for bills
CREATE TYPE IF NOT EXISTS bill_type AS ENUM ('electricity', 'water', 'gas');

CREATE TYPE IF NOT EXISTS bill_status_type AS ENUM ('unpaid', 'paid', 'overdue');

-- Create bills table
CREATE TABLE IF NOT EXISTS bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    name TEXT,
    bill_type bill_type NOT NULL,
    bill_id INTEGER UNIQUE NOT NULL, 
    amount INTEGER,
    status bill_status_type NOT NULL DEFAULT 'unpaid',
    paid_at TIMESTAMPTZ,
    due_date DATE NOT NULL,
    image_id UUID,
    apartment_id UUID NOT NULL,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    bill_id UUID NOT NULL,
    payer_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    payment_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (payer_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
