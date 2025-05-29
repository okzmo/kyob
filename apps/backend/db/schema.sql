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
-- Name: abilities; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.abilities AS ENUM (
    'ADMIN',
    'MANAGE_CHANNELS',
    'MANAGE_ROLES',
    'MANAGE_SERVER',
    'MANAGE_EXPRESSIONS',
    'CHANGE_NICKNAME',
    'MANAGE_NICKNAMES',
    'BAN',
    'KICK',
    'MUTE',
    'ATTACH_FILES',
    'MANAGE_MESSAGES'
);


--
-- Name: channel_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.channel_type AS ENUM (
    'voice',
    'textual',
    'dm',
    'groups'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channels (
    id character varying(20) NOT NULL,
    server_id character varying(20) NOT NULL,
    name character varying(255) NOT NULL,
    type public.channel_type NOT NULL,
    description text,
    users character varying(20)[],
    roles character varying(20)[],
    x integer NOT NULL,
    y integer NOT NULL,
    active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: facts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.facts (
    id character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    icon character varying(255),
    label character varying(255),
    value character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: friends; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.friends (
    id character varying(20) NOT NULL,
    user_id character varying(255) NOT NULL,
    friend_id character varying(255) NOT NULL,
    accepted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: invites; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.invites (
    id character varying(20) NOT NULL,
    server_id character varying(20) NOT NULL,
    invite_id character varying(255) NOT NULL,
    expire_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: links; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.links (
    id character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    label character varying(255),
    url character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.messages (
    id character varying(20) NOT NULL,
    author_id character varying(20) NOT NULL,
    server_id character varying(20) NOT NULL,
    channel_id character varying(20) NOT NULL,
    content jsonb NOT NULL,
    mentions_users character varying(20)[],
    mentions_channels character varying(20)[],
    attached text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id character varying(20) NOT NULL,
    server_id character varying(20) NOT NULL,
    name character varying(255) NOT NULL,
    color character varying(255) NOT NULL,
    description text,
    abilities text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: server_membership; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.server_membership (
    id character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    server_id character varying(20) NOT NULL,
    roles character varying(20)[],
    x integer NOT NULL,
    y integer NOT NULL,
    joined_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: servers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.servers (
    id character varying(20) NOT NULL,
    owner_id character varying(20) NOT NULL,
    name character varying(255) NOT NULL,
    avatar character varying(255),
    banner character varying(255),
    description text,
    private boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tokens (
    id character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    token text NOT NULL,
    type character varying(255) NOT NULL,
    expire_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id character varying(20) NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    avatar character varying(255),
    banner character varying(255),
    about text,
    gradient_top character varying(255),
    gradient_bottom character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- Name: facts facts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.facts
    ADD CONSTRAINT facts_pkey PRIMARY KEY (id);


--
-- Name: friends friends_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friends
    ADD CONSTRAINT friends_pkey PRIMARY KEY (id);


--
-- Name: invites invites_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invites
    ADD CONSTRAINT invites_pkey PRIMARY KEY (id);


--
-- Name: links links_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT links_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: server_membership server_membership_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.server_membership
    ADD CONSTRAINT server_membership_pkey PRIMARY KEY (id);


--
-- Name: server_membership server_membership_user_id_server_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.server_membership
    ADD CONSTRAINT server_membership_user_id_server_id_key UNIQUE (user_id, server_id);


--
-- Name: servers servers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.servers
    ADD CONSTRAINT servers_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_invites_invite_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_invites_invite_id ON public.invites USING btree (invite_id);


--
-- Name: idx_tokens_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tokens_token ON public.tokens USING btree (token);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: channels channels_server_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_server_id_fkey FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;


--
-- Name: facts facts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.facts
    ADD CONSTRAINT facts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: friends friends_friend_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friends
    ADD CONSTRAINT friends_friend_id_fkey FOREIGN KEY (friend_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: friends friends_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friends
    ADD CONSTRAINT friends_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: invites invites_server_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invites
    ADD CONSTRAINT invites_server_id_fkey FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;


--
-- Name: links links_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT links_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: messages messages_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: messages messages_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: messages messages_server_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_server_id_fkey FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;


--
-- Name: roles roles_server_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_server_id_fkey FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;


--
-- Name: server_membership server_membership_server_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.server_membership
    ADD CONSTRAINT server_membership_server_id_fkey FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;


--
-- Name: server_membership server_membership_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.server_membership
    ADD CONSTRAINT server_membership_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: servers servers_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.servers
    ADD CONSTRAINT servers_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: tokens tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250502134015');
