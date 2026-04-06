--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Ubuntu 14.13-1.pgdg20.04+1)
-- Dumped by pg_dump version 17.0

-- Started on 2026-01-29 10:24:17 WIB

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE IF EXISTS "ajinomoto-gate-system";
--
-- TOC entry 3381 (class 1262 OID 336075)
-- Name: ajinomoto-gate-system; Type: DATABASE; Schema: -; Owner: offonfarm
--

CREATE DATABASE "ajinomoto-gate-system" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE "ajinomoto-gate-system" OWNER TO offonfarm;

\connect -reuse-previous=on "dbname='ajinomoto-gate-system'"

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 214 (class 1259 OID 336132)
-- Name: app_config; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.app_config (
    id integer NOT NULL,
    app_name character varying NOT NULL,
    app_logo character varying,
    company_name character varying NOT NULL,
    company_email character varying,
    company_logo character varying,
    address character varying,
    smtp_host character varying,
    smtp_port integer,
    smtp_email character varying,
    smtp_password character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    updated_by uuid
);


ALTER TABLE public.app_config OWNER TO offonfarm;

--
-- TOC entry 213 (class 1259 OID 336131)
-- Name: app_config_id_seq; Type: SEQUENCE; Schema: public; Owner: offonfarm
--

CREATE SEQUENCE public.app_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.app_config_id_seq OWNER TO offonfarm;

--
-- TOC entry 3383 (class 0 OID 0)
-- Dependencies: 213
-- Name: app_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: offonfarm
--

ALTER SEQUENCE public.app_config_id_seq OWNED BY public.app_config.id;


--
-- TOC entry 209 (class 1259 OID 336079)
-- Name: auth_role; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.auth_role (
    id character varying(10) NOT NULL,
    name character varying(50),
    description text
);


ALTER TABLE public.auth_role OWNER TO offonfarm;

--
-- TOC entry 210 (class 1259 OID 336088)
-- Name: auth_user; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.auth_user (
    id character varying(36) NOT NULL,
    name character varying(255),
    username character varying(100) NOT NULL,
    email character varying(50) NOT NULL,
    password character varying(60) NOT NULL,
    role_id character varying(36) NOT NULL,
    person_id integer,
    status character varying(1),
    foto character varying,
    active boolean DEFAULT true,
    mobile_fcm_token character varying DEFAULT ''::character varying,
    web_fcm_token character varying DEFAULT ''::character varying,
    reset_token character varying,
    reset_token_expiry timestamp with time zone,
    created_by character varying(36),
    created_at timestamp with time zone,
    updated_by character varying(36),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.auth_user OWNER TO offonfarm;

--
-- TOC entry 211 (class 1259 OID 336104)
-- Name: c_menu; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.c_menu (
    id character varying(36) NOT NULL,
    name character varying(100),
    link character varying(100),
    icon character varying(20),
    description text,
    permission_label character varying(50),
    action character varying(50),
    level integer,
    seq integer,
    parent_id character varying(36),
    created_by character varying(36),
    created_at timestamp without time zone,
    updated_by character varying(36),
    updated_at timestamp without time zone,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.c_menu OWNER TO offonfarm;

--
-- TOC entry 212 (class 1259 OID 336112)
-- Name: c_menu_role; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.c_menu_role (
    id character varying(36) NOT NULL,
    menu_id character varying(36),
    role_id character varying(36),
    permission character varying(50),
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.c_menu_role OWNER TO offonfarm;

--
-- TOC entry 215 (class 1259 OID 336141)
-- Name: log_activity; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.log_activity (
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


ALTER TABLE public.log_activity OWNER TO offonfarm;

--
-- TOC entry 221 (class 1259 OID 336652)
-- Name: m_company_id_seq; Type: SEQUENCE; Schema: public; Owner: onlypkl
--

CREATE SEQUENCE public.m_company_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.m_company_id_seq OWNER TO onlypkl;

--
-- TOC entry 219 (class 1259 OID 336561)
-- Name: m_company; Type: TABLE; Schema: public; Owner: onlypkl
--

CREATE TABLE public.m_company (
    id integer DEFAULT nextval('public.m_company_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    is_registered_partner boolean DEFAULT false,
    pic_contact character varying(100),
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(36),
    updated_at timestamp with time zone,
    updated_by character varying(36),
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.m_company OWNER TO onlypkl;

--
-- TOC entry 222 (class 1259 OID 336653)
-- Name: m_department_id_seq; Type: SEQUENCE; Schema: public; Owner: onlypkl
--

CREATE SEQUENCE public.m_department_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.m_department_id_seq OWNER TO onlypkl;

--
-- TOC entry 220 (class 1259 OID 336575)
-- Name: m_department; Type: TABLE; Schema: public; Owner: onlypkl
--

CREATE TABLE public.m_department (
    id integer DEFAULT nextval('public.m_department_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    map_coordinate_json json,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(36),
    updated_at timestamp with time zone,
    updated_by character varying(36),
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.m_department OWNER TO onlypkl;

--
-- TOC entry 223 (class 1259 OID 336654)
-- Name: m_personel_id_seq; Type: SEQUENCE; Schema: public; Owner: onlypkl
--

CREATE SEQUENCE public.m_personel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.m_personel_id_seq OWNER TO onlypkl;

--
-- TOC entry 217 (class 1259 OID 336484)
-- Name: m_personnel; Type: TABLE; Schema: public; Owner: offonfarm
--

CREATE TABLE public.m_personnel (
    id integer DEFAULT nextval('public.m_personel_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    personnel_type character varying(50),
    departement_id integer,
    company_id integer,
    photo_face_template text,
    is_active boolean,
    created_at timestamp with time zone,
    created_by character varying(36),
    updated_at timestamp with time zone,
    updated_by character varying(36),
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.m_personnel OWNER TO offonfarm;

--
-- TOC entry 216 (class 1259 OID 336483)
-- Name: m_personnel_id_seq; Type: SEQUENCE; Schema: public; Owner: offonfarm
--

CREATE SEQUENCE public.m_personnel_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.m_personnel_id_seq OWNER TO offonfarm;

--
-- TOC entry 3384 (class 0 OID 0)
-- Dependencies: 216
-- Name: m_personnel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: offonfarm
--

ALTER SEQUENCE public.m_personnel_id_seq OWNED BY public.m_personnel.id;


--
-- TOC entry 218 (class 1259 OID 336538)
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: onlypkl
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO onlypkl;

--
-- TOC entry 3182 (class 2604 OID 336135)
-- Name: app_config id; Type: DEFAULT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.app_config ALTER COLUMN id SET DEFAULT nextval('public.app_config_id_seq'::regclass);


--
-- TOC entry 3366 (class 0 OID 336132)
-- Dependencies: 214
-- Data for Name: app_config; Type: TABLE DATA; Schema: public; Owner: offonfarm
--



--
-- TOC entry 3361 (class 0 OID 336079)
-- Dependencies: 209
-- Data for Name: auth_role; Type: TABLE DATA; Schema: public; Owner: offonfarm
--

INSERT INTO public.auth_role (id, name, description) VALUES ('HA01', 'SUPERADMIN', 'Hak Akses Paling Tinggi') ON CONFLICT DO NOTHING;
INSERT INTO public.auth_role (id, name, description) VALUES ('HA02', 'SECURITY', 'Hak Akses Security') ON CONFLICT DO NOTHING;


--
-- TOC entry 3362 (class 0 OID 336088)
-- Dependencies: 210
-- Data for Name: auth_user; Type: TABLE DATA; Schema: public; Owner: offonfarm
--

INSERT INTO public.auth_user (id, name, username, email, password, role_id, person_id, status, foto, active, mobile_fcm_token, web_fcm_token, reset_token, reset_token_expiry, created_by, created_at, updated_by, updated_at, deleted_at, is_deleted) VALUES ('972de936-15e4-48ed-ac85-1747de233fc5', 'Security', 'security', 'security@gmail.com', '$2a$10$RtfXxayvoIMWgRT8EoDU1OgQDd0p7TBA5eKsDS6rakJ8Ibc8vRZ9m', 'HA02', NULL, '1', NULL, true, '', '', NULL, NULL, '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 14:17:23.666149+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 14:17:30.346461+07', NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.auth_user (id, name, username, email, password, role_id, person_id, status, foto, active, mobile_fcm_token, web_fcm_token, reset_token, reset_token_expiry, created_by, created_at, updated_by, updated_at, deleted_at, is_deleted) VALUES ('019bb0fd-bb74-7edf-b097-7b677b481907', 'Superadmin', 'superadmin', 'superadmin@gmail.com', '$2a$10$vLixvpx4iTUYHRGcUQpdB.F0tmQ7Ha7IHJHIs2i5F6UvsIC8CHbRS', 'HA01', NULL, '1', NULL, true, '', '', NULL, NULL, '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 10:50:12.556364+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 14:15:41.416397+07', NULL, false) ON CONFLICT DO NOTHING;


--
-- TOC entry 3363 (class 0 OID 336104)
-- Dependencies: 211
-- Data for Name: c_menu; Type: TABLE DATA; Schema: public; Owner: offonfarm
--

INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('019bb0fd-bb74-7715-bbfd-fe91566cc00f', 'Dashboard', '/dashboard', 'ChartDotsIcon', 'Menu Dashboard', 'DASHBOARD', 'VIEW', 1, 1, NULL, NULL, '2026-01-12 13:59:54.374016', NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('5f0331d5-9181-48d4-a866-e768f6f54ff4', 'Master Data', '#', 'DatabaseIcon', 'Menu Master Data', 'MASTER_DATA', 'VIEW', 1, 2, NULL, '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:17:19.976424', NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('2a0dcaad-831e-4bd2-901c-dab55c6b75f5', 'Person', '/master/personnel', 'PointIcon', 'Master Data Personnel', 'PERSONNEL', 'VIEW,CREATE,UPDATE,DELETE', 2, 1, '5f0331d5-9181-48d4-a866-e768f6f54ff4', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:20:20.775983', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:20:32.971615', false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('12381af4-14ef-463d-86dd-0763ad4a7026', 'Departemen', '/master/departement', 'PointIcon', 'Master Data Departement', 'DEPARTEMENT', 'VIEW,CREATE,UPDATE,DELETE', 2, 2, '5f0331d5-9181-48d4-a866-e768f6f54ff4', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:21:53.998155', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:22:02.574202', false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('1a9a17de-8ca0-46a4-9d93-7982c2a3e5d1', 'Perusahaan', '/master/company', 'PointIcon', 'Master Data Perusahaan', 'COMPANY', 'VIEW,CREATE,UPDATE,DELETE', 2, 3, '5f0331d5-9181-48d4-a866-e768f6f54ff4', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:23:02.240302', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:23:09.84484', false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('1357f60b-aad1-41e3-80d5-f1f619ff5b03', 'Pengaturan', '#', 'SettingsIcon', 'Menu Pengaturan', 'SETTING', 'VIEW', 1, 3, NULL, '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:24:07.515264', NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('a4ee99fa-ff8a-4981-a3a9-632f2a2007d0', 'User', '/setting/user', 'PointIcon', 'Menu Pengaturan User', 'USER', 'VIEW,CREATE,UPDATE,DELETE', 2, 1, '1357f60b-aad1-41e3-80d5-f1f619ff5b03', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:25:27.926206', NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu (id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_by, created_at, updated_by, updated_at, is_deleted) VALUES ('dea4e4b9-deea-43a9-8cbc-adf4ca601f69', 'Menu', '/setting/menu', 'PointIcon', 'Menu Pengaturan Menu', 'MENU', 'VIEW,CREATE,DELETE,UPDATE', 2, 2, '1357f60b-aad1-41e3-80d5-f1f619ff5b03', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 13:26:11.264859', NULL, NULL, false) ON CONFLICT DO NOTHING;


--
-- TOC entry 3364 (class 0 OID 336112)
-- Dependencies: 212
-- Data for Name: c_menu_role; Type: TABLE DATA; Schema: public; Owner: offonfarm
--

INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('019bb0fd-bb74-74b3-ac45-1274aa52772b', '019bb0fd-bb74-7715-bbfd-fe91566cc00f', 'HA01', 'VIEW', '2026-01-27 13:16:13.084661') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('fdf59ca0-2cd1-4a91-9d4e-fd7c6598a015', '5f0331d5-9181-48d4-a866-e768f6f54ff4', 'HA01', 'VIEW', '2026-01-27 13:26:48.192215') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('81c6429e-69b8-44c2-9e8d-a76cb65db2d2', '2a0dcaad-831e-4bd2-901c-dab55c6b75f5', 'HA01', 'VIEW,CREATE,UPDATE,DELETE', '2026-01-27 13:26:48.19222') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('b5e05f90-434a-4938-82f3-e6be632c7e54', '12381af4-14ef-463d-86dd-0763ad4a7026', 'HA01', 'VIEW,CREATE,UPDATE,DELETE', '2026-01-27 13:26:48.192228') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('8f21f296-aebd-454b-9e0f-eacc78be95af', '1a9a17de-8ca0-46a4-9d93-7982c2a3e5d1', 'HA01', 'VIEW,CREATE,UPDATE,DELETE', '2026-01-27 13:26:48.192229') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('0708f8fa-7bac-487f-9982-0310f286289e', '1357f60b-aad1-41e3-80d5-f1f619ff5b03', 'HA01', 'VIEW', '2026-01-27 13:26:48.192236') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('69958f99-9891-4b00-925a-605478454337', 'a4ee99fa-ff8a-4981-a3a9-632f2a2007d0', 'HA01', 'VIEW,CREATE,UPDATE,DELETE', '2026-01-27 13:26:48.192237') ON CONFLICT DO NOTHING;
INSERT INTO public.c_menu_role (id, menu_id, role_id, permission, created_at) VALUES ('b3088846-c5cf-49ca-9bc7-142da0f3b76d', 'dea4e4b9-deea-43a9-8cbc-adf4ca601f69', 'HA01', 'VIEW,CREATE,UPDATE,DELETE', '2026-01-27 13:26:48.192238') ON CONFLICT DO NOTHING;


--
-- TOC entry 3367 (class 0 OID 336141)
-- Dependencies: 215
-- Data for Name: log_activity; Type: TABLE DATA; Schema: public; Owner: offonfarm
--

INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('259644ac-b866-4344-81b4-5c119c563749', 'Login', '2026-01-27 10:51:33.229853', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('57719673-7221-4e39-8f64-9d7c74995364', 'Login', '2026-01-27 10:52:06.773172', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('1fde037b-877a-47e0-aad5-ea3407e74717', 'Login', '2026-01-27 13:13:20.267293', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('d3437e66-f7ab-4375-9f93-fcd3b016ced3', 'Login', '2026-01-27 13:16:21.208642', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('3bad338c-9b25-4d48-8f97-922fd1f90f96', 'User', '2026-01-27 14:17:23.768374', 'Create User Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', '972de936-15e4-48ed-ac85-1747de233fc5') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('7555023f-567b-4062-86e5-453acc6a34dc', 'User', '2026-01-27 14:17:30.346483', 'Update User Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', '972de936-15e4-48ed-ac85-1747de233fc5') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('dad80c35-42c2-44d9-a2a9-3fc50682c13c', 'Login', '2026-01-27 14:45:08.587508', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('bcde8176-e54a-44fc-8098-0ef47a684632', 'Login', '2026-01-27 14:51:56.148144', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '', 'Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('a114a74f-219f-4e2f-8706-1572b5232291', 'Login', '2026-01-27 14:59:51.586624', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('ef0a1bea-701c-4bdf-8c75-73d339036e97', 'Login', '2026-01-27 15:16:16.925488', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '', 'Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('04cebf99-ad15-4453-8833-168e8dc9d4f1', 'Login', '2026-01-27 15:43:35.577435', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('cedd35b5-3387-4f1e-8ed0-d89a5652840d', 'Login', '2026-01-27 16:05:08.281026', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '', 'Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('cfcf30e5-d4be-4254-a096-3904badf7325', 'Login', '2026-01-27 16:34:20.864374', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('7ca2800c-4b0f-4f7e-bc24-adbad8eb4c42', 'Login', '2026-01-27 16:43:25.569064', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('d98eedd2-ba62-4d1f-aaa3-85d966512e9a', 'Login', '2026-01-27 16:58:22.458656', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('f18dc69f-9336-4f30-8112-12232f0048af', 'Login', '2026-01-27 17:31:33.026956', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('18e83cb0-3e6e-4519-997f-7f065480ea39', 'Login', '2026-01-27 17:50:07.099534', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('62230c7e-c84c-423b-b5f0-fcdbfe512d4c', 'Login', '2026-01-28 08:13:04.112843', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('cd8d87e1-1a41-4300-9d2b-0378525e534f', 'Login', '2026-01-28 08:17:13.407046', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('0ce743b4-f4fc-4e2c-aa10-d250894b492a', 'Login', '2026-01-28 08:22:11.634129', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('8ccf4e70-dc2b-4640-b0c6-75419a077290', 'Login', '2026-01-28 08:24:13.010464', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('9df57fa9-92ab-4db4-ae84-7655bb2dd355', 'Login', '2026-01-28 08:36:03.49023', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('c8ea155f-6f79-4312-8643-c2c941399fea', 'Login', '2026-01-28 08:41:07.890923', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('8b6104be-3c82-4a0a-9614-75a67a98fa74', 'Login', '2026-01-28 09:44:56.766354', 'Login Gagal Username Tidak Ditemukan', '00000000-0000-0000-0000-000000000000', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'asas') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('04f055bb-2a39-4201-bdb8-1b2d71cc92e2', 'Login', '2026-01-28 09:44:59.133317', 'Login Gagal Username Tidak Ditemukan', '00000000-0000-0000-0000-000000000000', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'asas') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('b1dc009a-7c38-4424-9888-e528f170a092', 'Login', '2026-01-28 09:44:59.305155', 'Login Gagal Username Tidak Ditemukan', '00000000-0000-0000-0000-000000000000', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'asas') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('1134e266-5fad-4753-9a5c-e6a8e8297fc1', 'Login', '2026-01-28 09:45:34.742277', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('6ede2054-5378-4a56-91b1-6d699603fcbc', 'Login', '2026-01-28 10:03:55.922127', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('a8ebf93b-dc61-48e3-a0fd-1e3f153e33eb', 'Login', '2026-01-28 10:15:32.848993', 'Login Gagal Username Tidak Ditemukan', '00000000-0000-0000-0000-000000000000', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmi') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('00be48e9-c3bf-440f-b4c5-460320036b34', 'Login', '2026-01-28 10:15:37.211643', 'Login Gagal Username Tidak Ditemukan', '00000000-0000-0000-0000-000000000000', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmi') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('1474b099-2d03-46af-8946-2f085ff2c403', 'Login', '2026-01-28 10:15:46.915329', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('483512b3-a540-44cf-b99d-9f454a40056a', 'Login', '2026-01-28 10:25:35.448759', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('3473ae5f-25d5-4c45-992e-0a58b33e8e64', 'Login', '2026-01-28 10:36:57.728047', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('dbec024b-8694-4ce7-b4f9-9ee6aff922ad', 'Login', '2026-01-28 10:48:39.363013', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;
INSERT INTO public.log_activity (id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) VALUES ('e21cdefe-7e93-4e87-8d3c-ec7da79d1c7c', 'Login', '2026-01-28 17:08:22.945581', 'Login Berhasil', '019bb0fd-bb74-7edf-b097-7b677b481907', 'WEB', '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36', 'superadmin') ON CONFLICT DO NOTHING;


--
-- TOC entry 3371 (class 0 OID 336561)
-- Dependencies: 219
-- Data for Name: m_company; Type: TABLE DATA; Schema: public; Owner: onlypkl
--

INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (1, 'PT Greatsoft Solusi Indonesia', true, '081234567890', '2026-01-27 14:59:16.997238+07', '019bb0fd-bb74-7edf-b097-7b677b481907', NULL, NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (2, 'PT Gajah Duduk', true, '0812', '2026-01-27 15:24:36.591209+07', '019bb0fd-bb74-7edf-b097-7b677b481907', NULL, NULL, NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (3, 'PT Kupu Kupu', true, '0811', '2026-01-27 15:24:46.576411+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 15:30:56.812029+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 15:30:56.812029+07', true) ON CONFLICT DO NOTHING;
INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (4, 'PT. ABC', false, 'John Doe', '2026-01-27 15:54:58.259768+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 15:56:47.774915+07', '019bb0fd-bb74-7edf-b097-7b677b481907', NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (5, 'PT. XYZ', true, '1212', '2026-01-27 15:54:58.259768+07', '019bb0fd-bb74-7edf-b097-7b677b481907', '2026-01-27 15:56:47.774915+07', '019bb0fd-bb74-7edf-b097-7b677b481907', NULL, false) ON CONFLICT DO NOTHING;
INSERT INTO public.m_company (id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted) VALUES (10, 'PT. well', true, NULL, '2026-01-27 15:56:47.774915+07', '019bb0fd-bb74-7edf-b097-7b677b481907', NULL, NULL, NULL, false) ON CONFLICT DO NOTHING;


--
-- TOC entry 3372 (class 0 OID 336575)
-- Dependencies: 220
-- Data for Name: m_department; Type: TABLE DATA; Schema: public; Owner: onlypkl
--



--
-- TOC entry 3370 (class 0 OID 336538)
-- Dependencies: 218
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: onlypkl
--

INSERT INTO public.schema_migrations (version, dirty) VALUES (20241220072834, false) ON CONFLICT DO NOTHING;


--
-- TOC entry 3385 (class 0 OID 0)
-- Dependencies: 213
-- Name: app_config_id_seq; Type: SEQUENCE SET; Schema: public; Owner: offonfarm
--

SELECT pg_catalog.setval('public.app_config_id_seq', 1, false);


--
-- TOC entry 3386 (class 0 OID 0)
-- Dependencies: 221
-- Name: m_company_id_seq; Type: SEQUENCE SET; Schema: public; Owner: onlypkl
--

SELECT pg_catalog.setval('public.m_company_id_seq', 10, true);


--
-- TOC entry 3387 (class 0 OID 0)
-- Dependencies: 222
-- Name: m_department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: onlypkl
--

SELECT pg_catalog.setval('public.m_department_id_seq', 1, false);


--
-- TOC entry 3388 (class 0 OID 0)
-- Dependencies: 223
-- Name: m_personel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: onlypkl
--

SELECT pg_catalog.setval('public.m_personel_id_seq', 1, false);


--
-- TOC entry 3389 (class 0 OID 0)
-- Dependencies: 216
-- Name: m_personnel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: offonfarm
--

SELECT pg_catalog.setval('public.m_personnel_id_seq', 1, false);


--
-- TOC entry 3204 (class 2606 OID 336140)
-- Name: app_config app_config_pkey; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.app_config
    ADD CONSTRAINT app_config_pkey PRIMARY KEY (id);


--
-- TOC entry 3194 (class 2606 OID 336087)
-- Name: auth_role auth_role_name_key; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.auth_role
    ADD CONSTRAINT auth_role_name_key UNIQUE (name);


--
-- TOC entry 3196 (class 2606 OID 336085)
-- Name: auth_role auth_role_pkey; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.auth_role
    ADD CONSTRAINT auth_role_pkey PRIMARY KEY (id);


--
-- TOC entry 3198 (class 2606 OID 336098)
-- Name: auth_user auth_user_pkey; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.auth_user
    ADD CONSTRAINT auth_user_pkey PRIMARY KEY (id);


--
-- TOC entry 3200 (class 2606 OID 336111)
-- Name: c_menu c_menu_pk; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.c_menu
    ADD CONSTRAINT c_menu_pk PRIMARY KEY (id);


--
-- TOC entry 3202 (class 2606 OID 336117)
-- Name: c_menu_role c_menu_role_pk; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.c_menu_role
    ADD CONSTRAINT c_menu_role_pk PRIMARY KEY (id);


--
-- TOC entry 3206 (class 2606 OID 336147)
-- Name: log_activity log_activity_pkey; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.log_activity
    ADD CONSTRAINT log_activity_pkey PRIMARY KEY (id);


--
-- TOC entry 3212 (class 2606 OID 336570)
-- Name: m_company m_company_name_key; Type: CONSTRAINT; Schema: public; Owner: onlypkl
--

ALTER TABLE ONLY public.m_company
    ADD CONSTRAINT m_company_name_key UNIQUE (name);


--
-- TOC entry 3214 (class 2606 OID 336568)
-- Name: m_company m_company_pkey; Type: CONSTRAINT; Schema: public; Owner: onlypkl
--

ALTER TABLE ONLY public.m_company
    ADD CONSTRAINT m_company_pkey PRIMARY KEY (id);


--
-- TOC entry 3216 (class 2606 OID 336585)
-- Name: m_department m_department_name_key; Type: CONSTRAINT; Schema: public; Owner: onlypkl
--

ALTER TABLE ONLY public.m_department
    ADD CONSTRAINT m_department_name_key UNIQUE (name);


--
-- TOC entry 3218 (class 2606 OID 336583)
-- Name: m_department m_department_pkey; Type: CONSTRAINT; Schema: public; Owner: onlypkl
--

ALTER TABLE ONLY public.m_department
    ADD CONSTRAINT m_department_pkey PRIMARY KEY (id);


--
-- TOC entry 3208 (class 2606 OID 336492)
-- Name: m_personnel m_personnel_pkey; Type: CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.m_personnel
    ADD CONSTRAINT m_personnel_pkey PRIMARY KEY (id);


--
-- TOC entry 3210 (class 2606 OID 336542)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: onlypkl
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 3219 (class 2606 OID 336099)
-- Name: auth_user auth_user_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.auth_user
    ADD CONSTRAINT auth_user_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.auth_role(id);


--
-- TOC entry 3220 (class 2606 OID 336118)
-- Name: c_menu_role c_menu_role_menu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.c_menu_role
    ADD CONSTRAINT c_menu_role_menu_id_fkey FOREIGN KEY (menu_id) REFERENCES public.c_menu(id);


--
-- TOC entry 3221 (class 2606 OID 336123)
-- Name: c_menu_role c_menu_role_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: offonfarm
--

ALTER TABLE ONLY public.c_menu_role
    ADD CONSTRAINT c_menu_role_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.auth_role(id);


--
-- TOC entry 3382 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2026-01-29 10:24:19 WIB

--
-- PostgreSQL database dump complete
--
