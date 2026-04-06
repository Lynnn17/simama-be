CREATE TABLE IF NOT EXISTS public.log_activity
(
    id character varying(36) NOT NULL,
    actions character varying(255),
    jam timestamp without time zone,
    keterangan text,
    id_user character varying(36),
    platform character varying(10),
    ip_address text,
    user_agent text,
    kode text
);
