-- remove timezone field from user table if exists
ALTER TABLE users DROP COLUMN IF EXISTS timezone;