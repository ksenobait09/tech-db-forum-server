CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS public.users
(
    nickname CITEXT NOT NULL,
    fullname varchar NOT NULL,
    email CITEXT NOT NULL,
    about varchar,
    CONSTRAINT users_pkey PRIMARY KEY (nickname)
);

CREATE UNIQUE INDEX IF NOT EXISTS index_users_email ON public.users(email);

CREATE TABLE IF NOT EXISTS public.forums
(
    posts integer NOT NULL DEFAULT 0,
    slug CITEXT NOT NULL,
    threads integer NOT NULL DEFAULT 0,
    title varchar NOT NULL,
    "user" CITEXT NOT NULL REFERENCES users(nickname),
    CONSTRAINT forums_pkey PRIMARY KEY (slug)
);

CREATE INDEX IF NOT EXISTS index_forums_user ON public.forums("user");

CREATE TABLE IF NOT EXISTS public.threads
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

CREATE INDEX IF NOT EXISTS index_threads_forum ON public.threads(forum);

CREATE TABLE IF NOT EXISTS public.votes
(
    idThread INT NOT NULL REFERENCES threads(id),
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    voice smallint NOT NULL,
    PRIMARY KEY (idThread, nickname)
);

CREATE INDEX IF NOT EXISTS index_votes_cover ON public.votes(idThread, nickname, voice);

CREATE TABLE IF NOT EXISTS public.posts (
    id SERIAL PRIMARY KEY,
    author CITEXT NOT NULL REFERENCES users(nickname),
    created TIMESTAMPTZ DEFAULT transaction_timestamp() NOT NULL,
    forum CITEXT NOT NULL,
    isEdited BOOLEAN NOT NULL DEFAULT FALSE,
    message VARCHAR NOT NULL,
    parent INTEGER DEFAULT 0 NOT NULL,
    path INTEGER[] DEFAULT array[]::INT[],
    rootParent INTEGER DEFAULT 0,
    thread INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS index_posts_thread_parent_id ON public.posts(thread, parent, id);
CREATE INDEX IF NOT EXISTS index_posts_thread_id ON public.posts(thread, id);
CREATE INDEX IF NOT EXISTS index_posts_thread_path ON public.posts(thread, path);
CREATE INDEX IF NOT EXISTS index_posts_rootparent_path ON public.posts(rootParent, path);

CREATE TABLE IF NOT EXISTS public.userforum (
    slug CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    PRIMARY KEY (slug, nickname)
)

