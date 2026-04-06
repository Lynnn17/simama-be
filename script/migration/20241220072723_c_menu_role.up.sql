CREATE TABLE IF NOT EXISTS public.c_menu_role
(
    id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    menu_id character varying(36) COLLATE pg_catalog."default",
    role_id character varying(36) COLLATE pg_catalog."default",
    permission character varying(50) COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT c_menu_role_pk PRIMARY KEY (id),
    CONSTRAINT c_menu_role_menu_id_fkey FOREIGN KEY (menu_id)
        REFERENCES public.c_menu (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT c_menu_role_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.auth_role (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)