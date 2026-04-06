CREATE TABLE IF NOT EXISTS public.company
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    code character varying(20) COLLATE pg_catalog."default" NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    refer_code character varying(100) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone,
    created_by character varying(36) COLLATE pg_catalog."default",
    updated_at timestamp without time zone,
    updated_by character varying(36) COLLATE pg_catalog."default",
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,

    CONSTRAINT company_pkey PRIMARY KEY (id),
    CONSTRAINT company_code_key UNIQUE (code),
    CONSTRAINT company_name_key UNIQUE (name),
    CONSTRAINT company_refer_code_key UNIQUE (refer_code)
)