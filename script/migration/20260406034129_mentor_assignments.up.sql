CREATE TABLE IF NOT EXISTS public.mentor_assignments
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    mentor_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    student_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    assigned_by character varying(255) COLLATE pg_catalog."default" NOT NULL,
    assigned_at timestamp with time zone DEFAULT now() NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    CONSTRAINT mentor_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT mentor_assignments_mentor_id_fkey FOREIGN KEY (mentor_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT mentor_assignments_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT mentor_assignments_assigned_by_fkey FOREIGN KEY (assigned_by)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
