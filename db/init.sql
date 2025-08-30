CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIUQE,
    phone TEXT UNIUQE,
    upi_vpa TEXT,
    creted_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS groups(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    creted_by UUID REFERENCES users(id),
    creted_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role TEXT DEFAULT `member`
    added_at TIMESTAMPZ NOT NULL DEFAULT now()
    PRIMARY KEY (group_id, user_id)
)