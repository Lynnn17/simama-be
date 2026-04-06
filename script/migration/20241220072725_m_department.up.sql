CREATE SEQUENCE IF NOT EXISTS public.m_department_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.m_department
(
    id integer NOT NULL DEFAULT nextval('m_department_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    map_coordinate_json json,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(36) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    updated_by character varying(36) COLLATE pg_catalog."default",
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,

    CONSTRAINT m_department_pkey PRIMARY KEY (id),
    CONSTRAINT m_department_name_key UNIQUE (name)
);