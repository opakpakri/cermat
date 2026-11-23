-- 1. Tabel Users (Setiap baris adalah 1 Tenant / Kantor Pemasaran / Perusahaan)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    business_name VARCHAR(100) NOT NULL DEFAULT 'Cermat - Cellix CRM',
    business_mode VARCHAR(20) NOT NULL DEFAULT 'OFFICE'
        CHECK (business_mode IN ('PROPERTY', 'OFFICE')),
    role VARCHAR(20) NOT NULL DEFAULT 'USER'
        CHECK (role IN ('ADMIN', 'USER')),
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'USER';

-- 2. Tabel Kontak / Pelanggan & Leads (Terisolasi per user_id)
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(150),
    company_name VARCHAR(100),
    address TEXT,
    lead_source VARCHAR(50) DEFAULT 'WHATSAPP',
    status VARCHAR(20) NOT NULL DEFAULT 'LEAD'
        CHECK (status IN ('LEAD', 'PROSPECT', 'CUSTOMER', 'INACTIVE')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, phone)
);
CREATE INDEX IF NOT EXISTS idx_contacts_user_id ON contacts(user_id);

-- 3. Tabel Katalog: Unit Properti atau Produk/Jasa Kantor (Terisolasi per user_id)
CREATE TABLE IF NOT EXISTS catalog_items (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(50),
    name VARCHAR(150) NOT NULL,
    category VARCHAR(50),
    price NUMERIC(15, 2) NOT NULL CHECK (price >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'BOOKED', 'SOLD', 'INACTIVE')),
    specs JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_catalog_items_user_status ON catalog_items(user_id, status);

-- 4. Tabel Deals / Sales Pipeline (Terisolasi per user_id)
CREATE TABLE IF NOT EXISTS deals (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    deal_code VARCHAR(30) NOT NULL,
    title VARCHAR(150) NOT NULL,
    contact_id INT NOT NULL REFERENCES contacts(id) ON DELETE RESTRICT,
    catalog_item_id INT REFERENCES catalog_items(id) ON DELETE SET NULL,
    stage VARCHAR(30) NOT NULL DEFAULT 'NEW'
        CHECK (stage IN ('NEW', 'QUALIFIED', 'PROPOSAL_SURVEY', 'NEGOTIATION_BOOKING', 'WON', 'LOST')),
    deal_value NUMERIC(15, 2) NOT NULL DEFAULT 0,
    payment_plan VARCHAR(50),
    expected_close_date DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, deal_code)
);
CREATE INDEX IF NOT EXISTS idx_deals_user_stage ON deals(user_id, stage);

-- 5. Tabel Aktivitas & Pengingat Follow-up (Terisolasi per user_id)
CREATE TABLE IF NOT EXISTS deal_activities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    deal_id INT REFERENCES deals(id) ON DELETE CASCADE,
    contact_id INT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    activity_type VARCHAR(30) NOT NULL
        CHECK (activity_type IN ('CALL', 'WHATSAPP', 'MEETING', 'SITE_VISIT', 'NOTE')),
    title VARCHAR(150) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_deal_activities_user_schedule ON deal_activities(user_id, is_completed, scheduled_at);

-- 6. Tabel Tiket Layanan / Keluhan Pelanggan (Terisolasi per user_id)
CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ticket_code VARCHAR(30) NOT NULL,
    contact_id INT NOT NULL REFERENCES contacts(id) ON DELETE RESTRICT,
    deal_id INT REFERENCES deals(id) ON DELETE SET NULL,
    subject VARCHAR(150) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM'
        CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH')),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, ticket_code)
);
CREATE INDEX IF NOT EXISTS idx_tickets_user_status ON tickets(user_id, status);
