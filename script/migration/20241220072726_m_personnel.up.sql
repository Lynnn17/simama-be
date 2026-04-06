CREATE SEQUENCE IF NOT EXISTS public.m_personnel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.m_personnel
(
    id integer NOT NULL DEFAULT nextval('m_personnel_id_seq'::regclass),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    personel_type character varying(20) COLLATE pg_catalog."default",
    department_id integer,
    company_id integer,
    photo_face_template integer,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(36) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    updated_by character varying(36) COLLATE pg_catalog."default",
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,

    CONSTRAINT m_personnel_pkey PRIMARY KEY (id),
    CONSTRAINT m_personnel_department_id_fkey FOREIGN KEY (department_id)
        REFERENCES public.m_department (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT m_personnel_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES public.m_company (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);