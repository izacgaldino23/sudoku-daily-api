-- add duration field in solves, the value is in seconds, the default value is 0, the calculated value is completed_at - started_at
ALTER TABLE solves ADD COLUMN IF NOT EXISTS duration INT NOT NULL DEFAULT 0;

-- update duration field in solves as int seconds not interval
UPDATE solves SET duration = EXTRACT(EPOCH FROM completed_at - started_at)::int;


-- remove completed_at field in solves
ALTER TABLE solves DROP COLUMN IF EXISTS completed_at;