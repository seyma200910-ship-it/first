CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    last_name VARCHAR(128) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role VARCHAR(32) NOT NULL DEFAULT 'worker'
        CHECK (role IN ('worker', 'admin', 'manager')),
    status VARCHAR(32) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'fired', 'on_leave', 'blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);