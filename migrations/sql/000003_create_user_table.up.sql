-- Provider enum
CREATE TYPE auth_provider AS ENUM ('email', 'google');

-- Create user table
CREATE TABLE IF NOT EXISTS "user" (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	provider auth_provider NOT NULL DEFAULT 'email',
	provider_id VARCHAR(255),
	email_verified BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);