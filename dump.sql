--
-- PostgreSQL database dump
--

-- Dumped from database version 12.5 (Ubuntu 12.5-1.pgdg20.04+1)
-- Dumped by pg_dump version 12.5 (Ubuntu 12.5-1.pgdg20.04+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id integer NOT NULL,
    content text NOT NULL,
    tags text[] DEFAULT '{}'::text[] NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.articles_id_seq OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_id_seq OWNED BY public.articles.id;


--
-- Name: articles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.articles_id_seq'::regclass);


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.articles VALUES (70, 'new text', '{"New tag","Super tag"}', '2021-01-07 09:04:02.13849');
INSERT INTO public.articles VALUES (90, '6666', '{''aaa'',''bob'',''cc''}', '2021-01-07 10:06:19.758585');
INSERT INTO public.articles VALUES (91, '8888', '{''aabob'',''acc''}', '2021-01-07 10:10:28.464973');


--
-- Name: articles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.articles_id_seq', 101, true);


--
-- Name: articles articles_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pk PRIMARY KEY (id);


--
-- Name: articles_tags_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX articles_tags_index ON public.articles USING gin (tags);


--
-- Name: articles on_insert; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER on_insert BEFORE INSERT ON public.articles FOR EACH ROW EXECUTE FUNCTION public.on_article_insert();


--
-- Name: articles on_update; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER on_update BEFORE UPDATE ON public.articles FOR EACH ROW EXECUTE FUNCTION public.on_article_update();


--
-- PostgreSQL database dump complete
--

