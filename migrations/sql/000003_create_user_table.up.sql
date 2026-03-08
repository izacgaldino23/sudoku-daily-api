-- Provider enum
DO $$
BEGIN
	IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'auth_provider') THEN
		CREATE TYPE auth_provider AS ENUM ('email', 'google');
	END IF;
END$$;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	provider auth_provider DEFAULT 'email',
	provider_id VARCHAR(255),
	email_verified BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);