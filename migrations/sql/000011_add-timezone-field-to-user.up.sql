-- Add timezone field to user table
ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone VARCHAR(255);