-- Optional: Enable foreign key constraints (run before anything else)
PRAGMA foreign_keys = ON;
-- Drop tables if they already exist
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS users_apartments;
DROP TABLE IF EXISTS apartment_invites;
DROP TABLE IF EXISTS apartments;
DROP TABLE IF EXISTS users;
-- USERS table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    -- UUID as TEXT
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- APARTMENTS table
CREATE TABLE IF NOT EXISTS apartments (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    unit_number INTEGER NOT NULL,
    admin_id TEXT NOT NULL,
    FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- USERS_APARTMENTS (many-to-many)
CREATE TABLE IF NOT EXISTS users_apartments (
    user_id TEXT NOT NULL,
    apartment_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    PRIMARY KEY (user_id, apartment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- APARTMENT_INVITES table
CREATE TABLE IF NOT EXISTS apartment_invites (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    invite_email TEXT NOT NULL,
    invite_status TEXT NOT NULL DEFAULT 'pending',
    -- ENUM emulated with CHECK
    invite_token TEXT UNIQUE NOT NULL,
    invite_expires_at DATETIME NOT NULL,
    apartment_id TEXT NOT NULL,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK (
        invite_status IN ('pending', 'accepted', 'declined')
    )
);
-- BILLS table
CREATE TABLE IF NOT EXISTS bills (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    name TEXT,
    bill_type TEXT NOT NULL,
    bill_id INTEGER UNIQUE NOT NULL,
    amount INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'unpaid',
    paid_at DATETIME,
    due_date DATE NOT NULL,
    image_id TEXT,
    apartment_id TEXT NOT NULL,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK (bill_type IN ('electricity', 'water', 'gas')),
    CHECK (status IN ('unpaid', 'paid', 'overdue'))
);
-- PAYMENTS table
CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    bill_id TEXT NOT NULL,
    payer_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    payment_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (payer_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);