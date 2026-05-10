ALTER TABLE public.logbooks ADD COLUMN IF NOT EXISTS reviewed_by character varying(36);
ALTER TABLE public.logbooks ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE public.logbooks ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE public.logbooks ALTER COLUMN progress_status SET DEFAULT 'pending';
ALTER TABLE public.logbooks RENAME COLUMN progress_status TO status;
