CREATE TABLE IF NOT EXISTS public.internship_registrations
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    university character varying(255) COLLATE pg_catalog."default" NOT NULL,
    major character varying(255) COLLATE pg_catalog."default" NOT NULL,
    cv_file_path character varying(500) COLLATE pg_catalog."default" NOT NULL,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT 'pending',
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT internship_registrations_pkey PRIMARY KEY (id),
    CONSTRAINT internship_registrations_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);