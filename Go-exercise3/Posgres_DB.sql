--
-- PostgreSQL database cluster dump
--

-- Started on 2024-04-16 17:05:14

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE postgres;
ALTER ROLE postgres WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'md5244af1e2823d5eaeeffc42c5096d8260';






--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

-- Dumped from database version 12.18
-- Dumped by pg_dump version 12.18

-- Started on 2024-04-16 17:05:14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Completed on 2024-04-16 17:05:14

--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

\connect postgres

--
-- PostgreSQL database dump
--

-- Dumped from database version 12.18
-- Dumped by pg_dump version 12.18

-- Started on 2024-04-16 17:05:14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 16384)
-- Name: adminpack; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS adminpack WITH SCHEMA pg_catalog;


--
-- TOC entry 2826 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION adminpack; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION adminpack IS 'administrative functions for PostgreSQL';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 204 (class 1259 OID 16395)
-- Name: passwords; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.passwords (
    id integer NOT NULL,
    password_string text NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    user_id integer NOT NULL
);


ALTER TABLE public.passwords OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 16393)
-- Name: passwords_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.passwords_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.passwords_id_seq OWNER TO postgres;

--
-- TOC entry 2827 (class 0 OID 0)
-- Dependencies: 203
-- Name: passwords_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.passwords_id_seq OWNED BY public.passwords.id;


--
-- TOC entry 2689 (class 2604 OID 16398)
-- Name: passwords id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.passwords ALTER COLUMN id SET DEFAULT nextval('public.passwords_id_seq'::regclass);


--
-- TOC entry 2820 (class 0 OID 16395)
-- Dependencies: 204
-- Data for Name: passwords; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.passwords (id, password_string, creation_date, user_id) FROM stdin;
1	generated_password	2024-04-16 15:27:16.646845	1
2	generated_password	2024-04-16 15:28:23.773375	1
3	generated_password	2024-04-16 15:30:59.376852	1
4	generated_password	2024-04-16 15:41:45.086365	1
5	generated_password	2024-04-16 15:41:47.349506	1
6	bvx~ePL7jw!E	2024-04-16 15:54:25.66212	1
7	]Cu7xipoPcAj	2024-04-16 15:54:50.911189	1
8	generated_password	2024-04-16 16:57:44.025712	1
9	!}]E6RDeJ4X(	2024-04-16 17:01:29.482517	1
10	nKns4Usy9@&x	2024-04-16 17:01:45.527319	1
11	xvM21U2rfr=+	2024-04-16 17:01:49.916448	1
\.


--
-- TOC entry 2828 (class 0 OID 0)
-- Dependencies: 203
-- Name: passwords_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.passwords_id_seq', 11, true);


--
-- TOC entry 2692 (class 2606 OID 16404)
-- Name: passwords passwords_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.passwords
    ADD CONSTRAINT passwords_pkey PRIMARY KEY (id);


-- Completed on 2024-04-16 17:05:15

--
-- PostgreSQL database dump complete
--

-- Completed on 2024-04-16 17:05:15

--
-- PostgreSQL database cluster dump complete
--

