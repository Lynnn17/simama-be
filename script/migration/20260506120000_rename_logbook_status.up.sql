ALTER TABLE public.logbooks RENAME COLUMN status TO progress_status;
ALTER TABLE public.logbooks ALTER COLUMN progress_status SET DEFAULT 'in_progress';
ALTER TABLE public.logbooks DROP COLUMN IF EXISTS notes;
ALTER TABLE public.logbooks DROP COLUMN IF EXISTS reviewed_at;
ALTER TABLE public.logbooks DROP COLUMN IF EXISTS reviewed_by;
