CREATE TABLE IF NOT EXISTS public.log_activity
(
    id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    actions character varying(255) COLLATE pg_catalog."default",
    jam timestamp without time zone,
    keterangan text COLLATE pg_catalog."default",
    id_user character varying(36) COLLATE pg_catalog."default",
    platform character varying(10) COLLATE pg_catalog."default",
    ip_address text COLLATE pg_catalog."default",
    user_agent text COLLATE pg_catalog."default",
    kode text COLLATE pg_catalog."default",
    CONSTRAINT log_activity_pkey PRIMARY KEY (id)
)
