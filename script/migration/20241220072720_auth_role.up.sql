CREATE TABLE IF NOT EXISTS public.auth_role
(
    id character varying(10) COLLATE pg_catalog."default" NOT NULL,
    name character varying(50) COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    CONSTRAINT auth_role_pkey PRIMARY KEY (id),
    CONSTRAINT auth_role_name_key UNIQUE (name)
)