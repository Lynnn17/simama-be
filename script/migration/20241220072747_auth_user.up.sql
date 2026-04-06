CREATE TABLE IF NOT EXISTS public.auth_user
(
    id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    name character varying(255) COLLATE pg_catalog."default",
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(60) COLLATE pg_catalog."default" NOT NULL,
    role_id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    status character varying(1) COLLATE pg_catalog."default",
    foto character varying COLLATE pg_catalog."default",
    active boolean DEFAULT true,
    mobile_fcm_token character varying COLLATE pg_catalog."default" DEFAULT ''::character varying,
    web_fcm_token character varying COLLATE pg_catalog."default" DEFAULT ''::character varying,
    reset_token character varying COLLATE pg_catalog."default",
    reset_token_expiry timestamp with time zone,
    created_by character varying(36) COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    updated_by character varying(36) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false,
    CONSTRAINT auth_user_pkey PRIMARY KEY (id),
    CONSTRAINT auth_user_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.auth_role (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
