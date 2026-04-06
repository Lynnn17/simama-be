CREATE SEQUENCE IF NOT EXISTS public.m_company_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.m_company
(
    id integer NOT NULL DEFAULT nextval('m_company_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    is_registered_partner boolean DEFAULT false,
    pic_contact character varying(100) COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(36) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    updated_by character varying(36) COLLATE pg_catalog."default",
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,

    CONSTRAINT m_company_pkey PRIMARY KEY (id),
    CONSTRAINT m_company_name_key UNIQUE (name)
)