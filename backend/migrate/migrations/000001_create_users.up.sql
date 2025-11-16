CREATE EXTENSION IF NOT exists citext;

CREATE TABLE user (
    id UUID PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    username CITEXT UNIQUE NOT NULL,
    prefix VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    provider VARCHAR(50),
    provider_id VARCHAR(255),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    password  bytea NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) WITH time zone NOT NULL DEFAULT now()
);