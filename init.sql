CREATE EXTENSION IF NOT EXISTS CITEXT;

-- Table: public.forums

-- DROP TABLE public.forums;

CREATE TABLE public.forums
(
    posts integer NOT NULL DEFAULT 0,
    slug varchar NOT NULL UNIQUE,
    threads integer NOT NULL DEFAULT 0,
    user_id integer NOT NULL,
    "user" CITEXT NOT NULL,
    CONSTRAINT forums_pkey PRIMARY KEY (slug)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.forums
    OWNER to docker;


-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE public.users
(
    nickname CITEXT NOT NULL,
    fullname varchar NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    about varchar,
    CONSTRAINT users_pkey PRIMARY KEY (nickname)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to docker;

/* TODO: finish thread */
CREATE TABLE public.thread
(
    author CITEXT NOT NULL,
    created TIMESTAMPTZ DEFAULT transaction_timestamp() NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (nickname)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.forums
    OWNER to docker;


CREATE TABLE public.post (
    id SERIAL PRIMARY KEY,
    author CITEXT NOT NULL,
    created TIMESTAMPTZ DEFAULT transaction_timestamp() NOT NULL,
    forum VARCHAR DEFAULT NULL,
    isEdited BOOLEAN DEFAULT FALSE,
    message varchar NOT NULL,
    parent INTEGER DEFAULT 0 NOT NULL,
    path varchar NULL,
    root INTEGER DEFAULT 0,
    thread INTEGER DEFAULT 0
);