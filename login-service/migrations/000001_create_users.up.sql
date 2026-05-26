-- 001_create_users.sql
-- Creates the users table for the auth service.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email        VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role         VARCHAR(50) NOT NULL DEFAULT 'customer',
    customer_id  VARCHAR(255) DEFAULT '',
    status       VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
