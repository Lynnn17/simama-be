CREATE TABLE IF NOT EXISTS public.section
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    code character varying(100) COLLATE pg_catalog."default" NOT NULL,
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    company_id uuid NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,

    CONSTRAINT section_pkey PRIMARY KEY (id),
    CONSTRAINT section_code_key UNIQUE (code),
    CONSTRAINT section_name_key UNIQUE (name),
    CONSTRAINT section_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES public.company (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)