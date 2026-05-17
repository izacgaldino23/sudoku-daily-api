CREATE INDEX idx_sudokus_size_date
ON public.sudokus (size, date);

CREATE INDEX idx_solves_duration 
ON public.solves (duration);

CREATE INDEX idx_solves_size 
ON public.solves (size);

CREATE INDEX idx_solves_user_started_at 
ON public.solves (user_id, started_at);