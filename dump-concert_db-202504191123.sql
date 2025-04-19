--
-- PostgreSQL database dump
--

-- Dumped from database version 15.12 (Debian 15.12-1.pgdg120+1)
-- Dumped by pg_dump version 16.8 (Ubuntu 16.8-0ubuntu0.24.04.1)

-- Started on 2025-04-19 11:23:26 WIB

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
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 3367 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 218 (class 1255 OID 16423)
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: admin
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO admin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 16397)
-- Name: bookings; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bookings (
    id integer NOT NULL,
    concert_id integer,
    user_id integer,
    quantity integer,
    booked_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.bookings OWNER TO admin;

--
-- TOC entry 216 (class 1259 OID 16396)
-- Name: bookings_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.bookings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.bookings_id_seq OWNER TO admin;

--
-- TOC entry 3368 (class 0 OID 0)
-- Dependencies: 216
-- Name: bookings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.bookings_id_seq OWNED BY public.bookings.id;


--
-- TOC entry 215 (class 1259 OID 16390)
-- Name: concerts; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.concerts (
    id integer NOT NULL,
    name_concert character varying(255),
    total_tickets integer,
    available_tickets integer,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.concerts OWNER TO admin;

--
-- TOC entry 214 (class 1259 OID 16389)
-- Name: concerts_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.concerts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.concerts_id_seq OWNER TO admin;

--
-- TOC entry 3369 (class 0 OID 0)
-- Dependencies: 214
-- Name: concerts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.concerts_id_seq OWNED BY public.concerts.id;


--
-- TOC entry 3208 (class 2604 OID 16400)
-- Name: bookings id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bookings ALTER COLUMN id SET DEFAULT nextval('public.bookings_id_seq'::regclass);


--
-- TOC entry 3205 (class 2604 OID 16393)
-- Name: concerts id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.concerts ALTER COLUMN id SET DEFAULT nextval('public.concerts_id_seq'::regclass);


--
-- TOC entry 3361 (class 0 OID 16397)
-- Dependencies: 217
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.bookings (id, concert_id, user_id, quantity, booked_at) FROM stdin;
\.


--
-- TOC entry 3359 (class 0 OID 16390)
-- Dependencies: 215
-- Data for Name: concerts; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.concerts (id, name_concert, total_tickets, available_tickets, start_time, end_time, created_at, updated_at) FROM stdin;
2	Justin beiber	1000	20	2025-04-19 01:05:30	2026-04-19 01:05:30	2025-04-19 11:09:13.639311	2025-04-19 11:11:54.503324
3	Justin beiber VIP	1000	10	2025-04-19 01:05:30	2025-04-30 01:05:30	2025-04-19 11:09:13.639311	2025-04-19 11:11:54.504792
1	coldplay	100	50	2025-04-15 01:05:30	2025-08-15 01:23:30	2025-04-19 07:22:41.329785	2025-04-19 11:12:05.330854
4	bruno mars	500	250	2025-04-19 01:05:30	2025-04-20 01:05:30	2025-04-19 07:23:55.668495	2025-04-19 11:12:05.394349
\.


--
-- TOC entry 3370 (class 0 OID 0)
-- Dependencies: 216
-- Name: bookings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.bookings_id_seq', 46, true);


--
-- TOC entry 3371 (class 0 OID 0)
-- Dependencies: 214
-- Name: concerts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.concerts_id_seq', 4, true);


--
-- TOC entry 3213 (class 2606 OID 16403)
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (id);


--
-- TOC entry 3211 (class 2606 OID 16395)
-- Name: concerts concerts_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.concerts
    ADD CONSTRAINT concerts_pkey PRIMARY KEY (id);


--
-- TOC entry 3215 (class 2620 OID 16424)
-- Name: concerts set_updated_at; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER set_updated_at BEFORE UPDATE ON public.concerts FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3214 (class 2606 OID 16404)
-- Name: bookings bookings_concert_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_concert_id_fkey FOREIGN KEY (concert_id) REFERENCES public.concerts(id);


-- Completed on 2025-04-19 11:23:27 WIB

--
-- PostgreSQL database dump complete
--

