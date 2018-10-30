CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE public.users
(
    nickname CITEXT NOT NULL,
    fullname varchar NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    about varchar,
    CONSTRAINT users_pkey PRIMARY KEY (nickname)
);

CREATE TABLE public.forums
(
    posts integer NOT NULL DEFAULT 0,
    slug CITEXT NOT NULL,
    threads integer NOT NULL DEFAULT 0,
    title varchar NOT NULL,
    "user" CITEXT NOT NULL REFERENCES users(nickname),
    CONSTRAINT forums_pkey PRIMARY KEY (slug)
);

CREATE TABLE public.threads
(
    id SERIAL PRIMARY KEY,
    author CITEXT NOT NULL REFERENCES users(nickname),
    created TIMESTAMPTZ,
    forum CITEXT NOT NULL REFERENCES forums(slug),
    message varchar NOT NULL,
    slug CITEXT DEFAULT NULL UNIQUE,
    title VARCHAR,
    votes INT NOT NULL DEFAULT 0
);

CREATE TABLE public.votes
(
    idThread INT NOT NULL REFERENCES threads(id),
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    voice smallint NOT NULL,
    PRIMARY KEY (idThread, nickname)
);

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