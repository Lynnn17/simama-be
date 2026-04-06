CREATE TABLE IF NOT EXISTS public.mentor_assignments
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    mentor_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    student_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT mentor_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT mentor_assignments_mentor_id_fkey FOREIGN KEY (mentor_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT mentor_assignments_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);