--
-- PostgreSQL database dump
--



DROP INDEX IF EXISTS _sys_announce_update_time_idx;
DROP INDEX IF EXISTS _sys_announce_type_idx;
DROP INDEX IF EXISTS _sys_announce_title_idx;
DROP INDEX IF EXISTS _sys_announce_status_idx;
DROP INDEX IF EXISTS _sys_announce_marked_idx;
ALTER TABLE IF EXISTS _sys_version_object ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS _sys_announce ALTER COLUMN tid DROP DEFAULT;
DROP SEQUENCE IF EXISTS _sys_version_object_tid_seq;
DROP TABLE IF EXISTS _sys_version_object;
DROP TABLE IF EXISTS _sys_object;
DROP TABLE IF EXISTS _sys_config;
DROP SEQUENCE IF EXISTS _sys_announce_tid_seq;
DROP TABLE IF EXISTS _sys_announce;


--
-- Name: _sys_announce; Type: TABLE; Schema: public;
--

CREATE TABLE _sys_announce (
    tid bigint NOT NULL,
    type integer NOT NULL,
    marked integer DEFAULT 0 NOT NULL,
    title character varying(1024) NOT NULL,
    info jsonb DEFAULT '{}'::jsonb NOT NULL,
    content jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN _sys_announce.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.tid IS 'the announce id';


--
-- Name: COLUMN _sys_announce.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.type IS 'the announce type, Normal=100:is normal type';


--
-- Name: COLUMN _sys_announce.marked; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.marked IS 'the announce marked';


--
-- Name: COLUMN _sys_announce.title; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.title IS 'the announce title';


--
-- Name: COLUMN _sys_announce.info; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.info IS 'the announce external info';


--
-- Name: COLUMN _sys_announce.content; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.content IS 'the announce content';


--
-- Name: COLUMN _sys_announce.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.update_time IS 'the announce update time';


--
-- Name: COLUMN _sys_announce.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.create_time IS 'the announce create time';


--
-- Name: COLUMN _sys_announce.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_announce.status IS 'the announce status, Normal=100:is normal status, Removed=-1:is removed status';


--
-- Name: _sys_announce_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE _sys_announce_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: _sys_announce_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE _sys_announce_tid_seq OWNED BY _sys_announce.tid;


--
-- Name: _sys_config; Type: TABLE; Schema: public;
--

CREATE TABLE _sys_config (
    key character varying(255) NOT NULL,
    value text NOT NULL,
    update_time timestamp with time zone NOT NULL
);


--
-- Name: _sys_object; Type: TABLE; Schema: public;
--

CREATE TABLE _sys_object (
    key character varying(255) NOT NULL,
    value json,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status smallint NOT NULL
);


--
-- Name: COLUMN _sys_object.key; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_object.key IS 'the object key';


--
-- Name: COLUMN _sys_object.value; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_object.value IS 'the object value';


--
-- Name: COLUMN _sys_object.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_object.create_time IS 'the create time';


--
-- Name: COLUMN _sys_object.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_object.status IS 'the status, Normal=100:is normal, Removed=-1:is removed';


--
-- Name: _sys_version_object; Type: TABLE; Schema: public;
--

CREATE TABLE _sys_version_object (
    tid bigint NOT NULL,
    key character varying(255) NOT NULL,
    pub text DEFAULT '*'::text NOT NULL,
    value jsonb,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status smallint NOT NULL
);


--
-- Name: COLUMN _sys_version_object.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.tid IS 'the primary key';


--
-- Name: COLUMN _sys_version_object.key; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.key IS 'the name of key';


--
-- Name: COLUMN _sys_version_object.pub; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.pub IS 'the publish scoe of version object, split multi by comma, * to all, x.x.x.x for ip';


--
-- Name: COLUMN _sys_version_object.value; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.value IS 'the version of key';


--
-- Name: COLUMN _sys_version_object.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.update_time IS 'the update time';


--
-- Name: COLUMN _sys_version_object.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.create_time IS 'the create time';


--
-- Name: COLUMN _sys_version_object.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN _sys_version_object.status IS 'the status, Normal=100:is normal, Disabled=200:is disabled, Removed=-1:is removed';


--
-- Name: _sys_version_object_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE _sys_version_object_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: _sys_version_object_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE _sys_version_object_tid_seq OWNED BY _sys_version_object.tid;


--
-- Name: _sys_announce tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_announce ALTER COLUMN tid SET DEFAULT nextval('_sys_announce_tid_seq'::regclass);


--
-- Name: _sys_version_object tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_version_object ALTER COLUMN tid SET DEFAULT nextval('_sys_version_object_tid_seq'::regclass);


--
-- Name: _sys_announce _sys_announce_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_announce
    ADD CONSTRAINT _sys_announce_pkey PRIMARY KEY (tid);


--
-- Name: _sys_config _sys_config_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_config
    ADD CONSTRAINT _sys_config_pkey PRIMARY KEY (key);


--
-- Name: _sys_object _sys_object_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_object
    ADD CONSTRAINT _sys_object_pkey PRIMARY KEY (key);


--
-- Name: _sys_version_object _sys_version_object_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY _sys_version_object
    ADD CONSTRAINT _sys_version_object_pkey PRIMARY KEY (tid);


--
-- Name: _sys_announce_marked_idx; Type: INDEX; Schema: public;
--

CREATE INDEX _sys_announce_marked_idx ON _sys_announce USING btree (marked);


--
-- Name: _sys_announce_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX _sys_announce_status_idx ON _sys_announce USING btree (status);


--
-- Name: _sys_announce_title_idx; Type: INDEX; Schema: public;
--

CREATE INDEX _sys_announce_title_idx ON _sys_announce USING btree (title);


--
-- Name: _sys_announce_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX _sys_announce_type_idx ON _sys_announce USING btree (type);


--
-- Name: _sys_announce_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX _sys_announce_update_time_idx ON _sys_announce USING btree (update_time);


--
-- PostgreSQL database dump complete
--

