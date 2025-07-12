BEGIN;

-- Enable pgcrypto extension (needed only once)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Safely create ENUM types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'invite_status_type') THEN
        CREATE TYPE invite_status_type AS ENUM ('pending', 'accepted', 'declined');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bill_type') THEN
        CREATE TYPE bill_type AS ENUM ('electricity', 'water', 'gas');
    END IF;
END$$;

-- USERS TABLE
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    first_name TEXT,
    last_name TEXT,
    email TEXT NOT NULL,
    password TEXT NOT NULL 
);

-- Add constraints and indexes
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_email'
    ) THEN
        CREATE UNIQUE INDEX idx_users_email ON users(email);
    END IF;

    -- Enforce email format (basic regex)
    -- BEGIN
    --     ALTER TABLE users ADD CONSTRAINT email_format CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$');
    -- EXCEPTION
    --     WHEN duplicate_object THEN NULL;
    -- END;

    -- BEGIN
    --     ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
    -- EXCEPTION
    --     WHEN duplicate_column THEN NULL;
    -- END;

    -- BEGIN
    --     ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;
    -- EXCEPTION
    --     WHEN duplicate_column THEN NULL;
    -- END;
END$$;

-- APARTMENTS TABLE
CREATE TABLE IF NOT EXISTS apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    unit_number INTEGER NOT NULL,
    admin_id UUID NOT NULL
);

-- Add foreign key for admin_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_apartments_admin_id'
    ) THEN
        ALTER TABLE apartments
        ADD CONSTRAINT fk_apartments_admin_id
        FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END$$;

-- USERS_APARTMENTS TABLE
CREATE TABLE IF NOT EXISTS users_apartments (
    user_id UUID NOT NULL,
    apartment_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    invite_status invite_status_type NOT NULL DEFAULT 'pending',
    invite_token TEXT UNIQUE NOT NULL,
    invite_expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, apartment_id)
);

-- Add foreign keys for users_apartments
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_users_apartments_user'
    ) THEN
        ALTER TABLE users_apartments
        ADD CONSTRAINT fk_users_apartments_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_users_apartments_apartment'
    ) THEN
        ALTER TABLE users_apartments
        ADD CONSTRAINT fk_users_apartments_apartment
        FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE;
    END IF;

    -- Optional: ensure invite hasn't expired
    BEGIN
        ALTER TABLE users_apartments
        ADD CONSTRAINT chk_invite_expiry CHECK (invite_expires_at > now());
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END;
END$$;

-- BILLS TABLE
CREATE TABLE IF NOT EXISTS bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    bill_type bill_type NOT NULL,
    bill_id INTEGER UNIQUE NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    due_date DATE NOT NULL,
    image_id UUID,
    apartment_id UUID NOT NULL
);

-- Add foreign key for bills
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_bills_apartment'
    ) THEN
        ALTER TABLE bills
        ADD CONSTRAINT fk_bills_apartment
        FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE;
    END IF;
END$$;

-- Add TRIGGER FUNCTION for auto-updating updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to each table (if not already exists)
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOREACH tbl IN ARRAY ['users', 'apartments', 'users_apartments', 'bills'] LOOP
        IF NOT EXISTS (
            SELECT 1 FROM pg_trigger WHERE tgname = format('trigger_update_%s', tbl)
        ) THEN
            EXECUTE format(
                'CREATE TRIGGER trigger_update_%1$I
                 BEFORE UPDATE ON %1$I
                 FOR EACH ROW
                 EXECUTE FUNCTION set_updated_at();',
                tbl
            );
        END IF;
    END LOOP;
END$$;

COMMIT;
