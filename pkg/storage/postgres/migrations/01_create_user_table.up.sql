CREATE TABLE IF NOT EXISTS users (
    user_id UUID NOT NULL PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    email_id VARCHAR(320) NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    user_name VARCHAR NOT NULL
);