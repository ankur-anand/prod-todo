CREATE TABLE IF NOT EXISTS users (
    user_id UUID NOT NULL PRIMARY KEY,
    email_id VARCHAR(320) NOT NULL UNIQUE,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);