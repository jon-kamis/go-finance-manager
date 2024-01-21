--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

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
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(255),
    create_dt timestamp without time zone,
    last_update_dt timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    code character varying(255) NOT NULL,
    create_dt timestamp without time zone,
    last_update_dt timestamp without time zone
);


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.roles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    id integer NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    code character varying(255) NOT NULL,
    create_dt timestamp without time zone,
    last_update_dt timestamp without time zone
);


--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.user_roles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: bank_account; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bank_accounts (
    id integer NOT NULL,
    account_name character varying(255) NOT NULL,
    bank_name character varying(255) NOT NULL,
    balance NUMERIC(10,2) NOT NULL,
    user_id integer NOT NULL,
    create_dt timestamp without time zone,
    last_update_dt timestamp without time zone
);

--
-- Name: bank_account_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.bank_accounts ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.bank_account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: loans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loans (
    id integer NOT NULL,
    user_id integer NOT NULL,
    loan_name character varying(255) NOT NULL,
    total_balance NUMERIC(10,2) NOT NULL,
    total_cost NUMERIC(10,2) NOT NULL,
    total_principal NUMERIC(10,2) NOT NULL,
    total_interest NUMERIC(10,2),
    monthly_payment NUMERIC(10,2),
    interest_rate NUMERIC(10,5),
    loan_term integer not null,
    create_dt timestamp without time zone,
    last_update_dt timestamp without time zone
);

--
-- Name: loan_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.loans ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.loan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: incomes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.incomes (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    type character varying(255) NOT NULL,
    rate NUMERIC(10, 2) NOT NULL,
    hours NUMERIC(10, 2) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    frequency character varying(255) NOT NULL,
    tax_percentage NUMERIC(10, 2) NOT NULL,
    start_dt timestamp,
    create_dt timestamp,
    last_update_dt timestamp without time zone
);

--
-- Name: incomes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.incomes ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.income_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: stocks; Type: Table; Schema: public; Owner: -
--
CREATE TABLE public.stocks (
    id integer NOT NULL,
    ticker character varying(255) NOT NULL,
    high NUMERIC(10, 4) NOT NULL,
    low NUMERIC(10, 4) NOT NULL,
    open NUMERIC(10, 4) NOT NULL,
    close NUMERIC(10, 4) NOT NULL,
    date timestamp,
    create_dt timestamp,
    last_update_dt timestamp
);

--
-- Name: stocks_id_seq; Type: SEQUENCE; Schema: public; Owner -
--
ALTER TABLE public.stocks ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: stock_data; Type: Table; Schema: public; Owner: -
--
CREATE TABLE public.stock_data (
    id integer NOT NULL,
    ticker character varying(255) NOT NULL,
    high NUMERIC(10, 4) NOT NULL,
    low NUMERIC(10, 4) NOT NULL,
    open NUMERIC(10, 4) NOT NULL,
    close NUMERIC(10, 4) NOT NULL,
    date timestamp,
    create_dt timestamp,
    last_update_dt timestamp
);

--
-- Name: stock_data_id_seq; Type: SEQUENCE; Schema: public; Owner -
--
ALTER TABLE public.stock_data ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.stock_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: user_stocks; Type: Table; Schema: public; Owner: -
--
CREATE TABLE public.user_stocks (
    id integer NOT NULL,
    user_id integer NOT NULL,
    ticker character varying(255) NOT NULL,
    quantity NUMERIC(10,4) NOT NULL,
    effective_dt timestamp NOT NULL,
    expiration_dt timestamp NOT NULL,
    create_dt timestamp,
    last_update_dt timestamp
);

--
-- Name: user_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner -
--
ALTER TABLE public.user_stocks ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: bills; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bills (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    create_dt timestamp,
    last_update_dt timestamp
);

--
-- Name: bills_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.bills ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.bill_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: credit_cards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.credit_cards (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    balance NUMERIC(10, 2) NOT NULL,
    credit_limit NUMERIC(10, 2) NOT NULL,
    apr NUMERIC(10, 2) NOT NULL,
    min_pay NUMERIC(10, 2) NOT NULL,
    min_pay_percentage NUMERIC(10, 2) NOT NULL,
    create_dt timestamp,
    last_update_dt timestamp
);

--
-- Name: bills_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.credit_cards ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.credit_card_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

COPY public.users (id, username, first_name, last_name, email, password, create_dt, last_update_dt) FROM stdin;
1	admin	admin	istrator	admin@fm.com	$2a$10$S9nLk.BzkZuSPXvdn6JXoO0VX/tf8QNebc0ct8J39n.mU8Gzz.pPS	2023-11-13 00:00:00	2023-11-13 00:00:00
\.

SELECT pg_catalog.setval('public.users_id_seq', 2, true);

COPY public.roles (id, code, create_dt, last_update_dt) FROM stdin;
1	admin	2023-11-13 00:00:00	2023-11-13 00:00:00
2	user	2023-11-13 00:00:00	2023-11-13 00:00:00
\.

SELECT pg_catalog.setval('public.roles_id_seq', 3, true);

COPY public.user_roles (id, user_id, role_id, code, create_dt, last_update_dt) FROM stdin;
1	1	1	admin	2023-11-13 00:00:00	2023-11-13 00:00:00
2	1	2	user	2023-11-13 00:00:00	2023-11-13 00:00:00
\.

SELECT pg_catalog.setval('public.user_roles_id_seq', 3, true);

--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_pkey PRIMARY KEY (id);

--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

--
-- Name: bank_accounts bank_account_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bank_accounts
    ADD CONSTRAINT bank_account_pkey PRIMARY KEY (id);

--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);

--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);

--
-- PostgreSQL database dump complete
--