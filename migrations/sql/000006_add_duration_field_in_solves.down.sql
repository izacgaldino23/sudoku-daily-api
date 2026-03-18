-- add completed_at field in solves
ALTER TABLE solves ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- update completed_at field in solves
UPDATE solves SET completed_at = (started_at + duration);

-- remove duration field in solves
ALTER TABLE solves DROP COLUMN IF EXISTS duration;