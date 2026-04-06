CREATE TABLE IF NOT EXISTS public.logbooks
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    student_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    log_date date NOT NULL,
    activities text NOT NULL,
    blockers text NOT NULL,
    plan_tomorrow text NOT NULL,
    evidence_url character varying COLLATE pg_catalog."default",
    status character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT 'pending',
    notes text,
    submitted_at timestamp with time zone DEFAULT now(),
    reviewed_at timestamp with time zone,
    reviewed_by character varying(36) COLLATE pg_catalog."default",
    CONSTRAINT logbooks_pkey PRIMARY KEY (id),
    CONSTRAINT logbooks_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT logbooks_reviewed_by_fkey FOREIGN KEY (reviewed_by)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
);