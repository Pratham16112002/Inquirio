

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE IF NOT EXISTS mentor (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username CITEXT UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    provider VARCHAR(50), 
    provider_id VARCHAR(255),  
    password BYTEA,            
    email CITEXT UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    experience_years NUMERIC(5,2),
    bio TEXT,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) WITH time zone NOT NULL DEFAULT now()
);