-- Enable UUID generator extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create enum types if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'invite_status_type') THEN
        CREATE TYPE invite_status_type AS ENUM ('pending', 'accepted', 'declined');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bill_type') THEN
        CREATE TYPE bill_type AS ENUM ('electricity', 'water', 'gas');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bill_status_type') THEN
        CREATE TYPE bill_status_type AS ENUM ('unpaid', 'paid', 'overdue');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status_type') THEN
        CREATE TYPE payment_status_type AS ENUM ('pending', 'paid', 'failed', 'cancelled');
    END IF;
END $$;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Apartments table
CREATE TABLE IF NOT EXISTS apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    unit_number INTEGER NOT NULL,
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Users ↔ Apartments junction table
CREATE TABLE IF NOT EXISTS users_apartments (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    apartment_id UUID NOT NULL REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (user_id, apartment_id),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- Apartment invites
CREATE TABLE IF NOT EXISTS apartment_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ,
    invite_email TEXT NOT NULL,
    invite_status invite_status_type NOT NULL DEFAULT 'pending',
    invite_token TEXT UNIQUE NOT NULL,
    invite_expires_at TIMESTAMPTZ NOT NULL,
    apartment_id UUID NOT NULL REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Bills table
CREATE TABLE IF NOT EXISTS bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    bill_type bill_type NOT NULL,
    bill_id INTEGER UNIQUE NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    -- status bill_status_type NOT NULL DEFAULT 'unpaid',
    -- paid_at TIMESTAMPTZ,
    due_date DATE NOT NULL,
    image_id UUID,
    apartment_id UUID NOT NULL REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ,
    bill_id UUID NOT NULL REFERENCES bills(id) ON DELETE CASCADE ON UPDATE CASCADE,
    payer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    paid_at TIMESTAMPTZ,
    status payment_status_type NOT NULL DEFAULT 'pending',
    gateway TEXT NOT NULL,
    transaction_id TEXT,
    callback_data JSONB
);
