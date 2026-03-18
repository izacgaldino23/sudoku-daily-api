-- add size field in solves NOT NULL CHECK (size in (4, 6, 9))
ALTER TABLE solves ADD COLUMN IF NOT EXISTS size INT NOT NULL CHECK (size in (4, 6, 9));