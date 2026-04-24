CREATE TABLE IF NOT EXISTS public.internship_registrations
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id character varying(255) COLLATE pg_catalog."default",
    reviewed_by character varying(255) COLLATE pg_catalog."default",
    full_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    university character varying(255) COLLATE pg_catalog."default" NOT NULL,
    major character varying(255) COLLATE pg_catalog."default" NOT NULL,
    semester character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(25) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    period character varying(50) COLLATE pg_catalog."default" NOT NULL,
    cv_file_path character varying(500) COLLATE pg_catalog."default" NOT NULL,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT 'pending',
    applied_at timestamp with time zone DEFAULT now() NOT NULL,
    reviewed_at timestamp with time zone,
    CONSTRAINT internship_registrations_pkey PRIMARY KEY (id),
    CONSTRAINT internship_registrations_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL,
    CONSTRAINT internship_registrations_reviewed_by_fkey FOREIGN KEY (reviewed_by)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
);
