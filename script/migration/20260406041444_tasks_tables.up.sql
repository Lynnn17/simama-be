CREATE TABLE IF NOT EXISTS public.tasks
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    mentor_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    student_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    title character varying(255) COLLATE pg_catalog."default" NOT NULL,
    description text NOT NULL,
    deadline timestamp with time zone NOT NULL,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT 'assigned',
    grade integer,
    feedback text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_mentor_id_fkey FOREIGN KEY (mentor_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT tasks_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS public.task_files
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    task_id uuid NOT NULL,
    file_url character varying COLLATE pg_catalog."default" NOT NULL,
    uploaded_by character varying(36) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT task_files_pkey PRIMARY KEY (id),
    CONSTRAINT task_files_task_id_fkey FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT task_files_uploaded_by_fkey FOREIGN KEY (uploaded_by)
        REFERENCES public.auth_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);