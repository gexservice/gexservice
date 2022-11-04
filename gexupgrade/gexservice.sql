--
-- PostgreSQL database dump
--

-- Dumped from database version 13.8 (Debian 13.8-1.pgdg110+1)
-- Dumped by pg_dump version 13.8 (Debian 13.8-1.pgdg110+1)


DROP INDEX IF EXISTS exs_withdraw_user_id_idx;
DROP INDEX IF EXISTS exs_withdraw_update_time_idx;
DROP INDEX IF EXISTS exs_withdraw_type_idx;
DROP INDEX IF EXISTS exs_withdraw_status_idx;
DROP INDEX IF EXISTS exs_withdraw_order_id_idx;
DROP INDEX IF EXISTS exs_withdraw_asset_idx;
DROP INDEX IF EXISTS exs_user_update_time_idx;
DROP INDEX IF EXISTS exs_user_type_idx;
DROP INDEX IF EXISTS exs_user_status_idx;
DROP INDEX IF EXISTS exs_user_role_idx;
DROP INDEX IF EXISTS exs_user_phone_idx;
DROP INDEX IF EXISTS exs_user_password_idx;
DROP INDEX IF EXISTS exs_user_account_idx;
DROP INDEX IF EXISTS exs_order_user_id_idx;
DROP INDEX IF EXISTS exs_order_update_time_idx;
DROP INDEX IF EXISTS exs_order_unhedged_idx;
DROP INDEX IF EXISTS exs_order_type_idx;
DROP INDEX IF EXISTS exs_order_trigger_idx;
DROP INDEX IF EXISTS exs_order_symobl_idx;
DROP INDEX IF EXISTS exs_order_status_idx;
DROP INDEX IF EXISTS exs_order_side_idx;
DROP INDEX IF EXISTS exs_order_order_id_idx;
DROP INDEX IF EXISTS exs_order_fee_settled_idx;
DROP INDEX IF EXISTS exs_order_comm_user_type_idx;
DROP INDEX IF EXISTS exs_order_comm_status_idx;
DROP INDEX IF EXISTS exs_order_comm_create_time_idx;
DROP INDEX IF EXISTS exs_kline_symbol_idx;
DROP INDEX IF EXISTS exs_kline_start_time_idx;
DROP INDEX IF EXISTS exs_kline_interval_idx;
DROP INDEX IF EXISTS exs_holding_user_symbol_idx;
DROP INDEX IF EXISTS exs_holding_update_time_idx;
DROP INDEX IF EXISTS exs_holding_status_idx;
DROP INDEX IF EXISTS exs_holding_blowup_idx;
DROP INDEX IF EXISTS exs_holding_amount_idx;
DROP INDEX IF EXISTS exs_balance_user_area_asset_idx;
DROP INDEX IF EXISTS exs_balance_status_idx;
DROP INDEX IF EXISTS exs_balance_history_user_asset_idx;
DROP INDEX IF EXISTS exs_balance_history_status_idx;
ALTER TABLE IF EXISTS exs_user ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_order_comm ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_order ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_kline ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_holding ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_balance_history ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS exs_balance ALTER COLUMN tid DROP DEFAULT;
DROP TABLE IF EXISTS exs_withdraw;
DROP SEQUENCE IF EXISTS exs_user_tid_seq;
DROP TABLE IF EXISTS exs_user;
DROP SEQUENCE IF EXISTS exs_order_tid_seq;
DROP SEQUENCE IF EXISTS exs_order_comm_tid_seq;
DROP TABLE IF EXISTS exs_order_comm;
DROP TABLE IF EXISTS exs_order;
DROP SEQUENCE IF EXISTS exs_kline_tid_seq;
DROP TABLE IF EXISTS exs_kline;
DROP SEQUENCE IF EXISTS exs_holding_tid_seq;
DROP TABLE IF EXISTS exs_holding;
DROP SEQUENCE IF EXISTS exs_balance_tid_seq;
DROP SEQUENCE IF EXISTS exs_balance_record_tid_seq;
DROP TABLE IF EXISTS exs_balance_history;
DROP TABLE IF EXISTS exs_balance;


--
-- Name: exs_balance; Type: TABLE; Schema: public;
--

CREATE TABLE exs_balance (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    area integer NOT NULL,
    asset character varying(30) NOT NULL,
    free double precision DEFAULT 0 NOT NULL,
    locked double precision DEFAULT 0 NOT NULL,
    margin double precision DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_balance.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.tid IS 'the primary key';


--
-- Name: COLUMN exs_balance.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.user_id IS 'the balance user id';


--
-- Name: COLUMN exs_balance.area; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.area IS 'the balance area, Funds=100:is funds area, Spot=200:is spot area, Futures=300:is futures area';


--
-- Name: COLUMN exs_balance.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.asset IS 'the balance asset key';


--
-- Name: COLUMN exs_balance.free; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.free IS 'the balance free amount';


--
-- Name: COLUMN exs_balance.locked; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.locked IS 'the balance locked amount';


--
-- Name: COLUMN exs_balance.margin; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.margin IS 'the balance margin value';


--
-- Name: COLUMN exs_balance.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.update_time IS 'the balance last update time';


--
-- Name: COLUMN exs_balance.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.create_time IS 'the balance create time';


--
-- Name: COLUMN exs_balance.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance.status IS 'the balance status, Normal=100: is normal, Locked=200: is locked';


--
-- Name: exs_balance_history; Type: TABLE; Schema: public;
--

CREATE TABLE exs_balance_history (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    asset character varying(30) NOT NULL,
    valuation double precision DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_balance_history.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.tid IS 'the primary key';


--
-- Name: COLUMN exs_balance_history.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.user_id IS 'the balance user id';


--
-- Name: COLUMN exs_balance_history.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.asset IS 'the balance asset key';


--
-- Name: COLUMN exs_balance_history.valuation; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.valuation IS 'the balance valuation';


--
-- Name: COLUMN exs_balance_history.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.update_time IS 'the balance record update time';


--
-- Name: COLUMN exs_balance_history.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.create_time IS 'the balance record create time, is daily zero time';


--
-- Name: COLUMN exs_balance_history.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_balance_history.status IS 'the balance record status, Normal=100: is normal status';


--
-- Name: exs_balance_record_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_balance_record_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_balance_record_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_balance_record_tid_seq OWNED BY exs_balance_history.tid;


--
-- Name: exs_balance_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_balance_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_balance_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_balance_tid_seq OWNED BY exs_balance.tid;


--
-- Name: exs_holding; Type: TABLE; Schema: public;
--

CREATE TABLE exs_holding (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    symbol character varying(16) NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    open double precision DEFAULT 0 NOT NULL,
    blowup double precision DEFAULT 0 NOT NULL,
    lever integer DEFAULT 1 NOT NULL,
    margin_used double precision DEFAULT 0 NOT NULL,
    margin_added double precision DEFAULT 0 NOT NULL,
    update_time timestamp(6) with time zone NOT NULL,
    create_time timestamp(6) with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_holding.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.tid IS 'the primary key';


--
-- Name: COLUMN exs_holding.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.user_id IS 'the holding user id';


--
-- Name: COLUMN exs_holding.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.symbol IS 'the holding symbol';


--
-- Name: COLUMN exs_holding.amount; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.amount IS 'the holding amount';


--
-- Name: COLUMN exs_holding.open; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.open IS 'the holding open price';


--
-- Name: COLUMN exs_holding.blowup; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.blowup IS 'the holding blowup price';


--
-- Name: COLUMN exs_holding.lever; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.lever IS 'the holding lever';


--
-- Name: COLUMN exs_holding.margin_used; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.margin_used IS 'the holding margin used';


--
-- Name: COLUMN exs_holding.margin_added; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.margin_added IS 'the holding margin added';


--
-- Name: COLUMN exs_holding.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.update_time IS 'the holding last update time';


--
-- Name: COLUMN exs_holding.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.create_time IS 'the holding create time';


--
-- Name: COLUMN exs_holding.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_holding.status IS 'the holding status, Normal=100: is normal, Locked=200: is locked';


--
-- Name: exs_holding_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_holding_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_holding_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_holding_tid_seq OWNED BY exs_holding.tid;


--
-- Name: exs_kline; Type: TABLE; Schema: public;
--

CREATE TABLE exs_kline (
    tid bigint NOT NULL,
    symbol character varying(30) NOT NULL,
    interv character varying(8) NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    count bigint DEFAULT 0 NOT NULL,
    open double precision DEFAULT 0 NOT NULL,
    close double precision DEFAULT 0 NOT NULL,
    low double precision DEFAULT 0 NOT NULL,
    high double precision DEFAULT 0 NOT NULL,
    volume double precision NOT NULL,
    start_time timestamp with time zone NOT NULL,
    update_time timestamp with time zone NOT NULL
);


--
-- Name: COLUMN exs_kline.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.tid IS 'the primay key';


--
-- Name: COLUMN exs_kline.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.symbol IS 'the kline symbol';


--
-- Name: COLUMN exs_kline.interv; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.interv IS 'the kline interval key';


--
-- Name: COLUMN exs_kline.amount; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.amount IS 'the kline amount';


--
-- Name: COLUMN exs_kline.count; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.count IS 'the kline count';


--
-- Name: COLUMN exs_kline.open; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.open IS 'the kline open price';


--
-- Name: COLUMN exs_kline.close; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.close IS 'the kline close price';


--
-- Name: COLUMN exs_kline.low; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.low IS 'the kline low price';


--
-- Name: COLUMN exs_kline.high; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.high IS 'the kline high price';


--
-- Name: COLUMN exs_kline.volume; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.volume IS 'the kline volume price';


--
-- Name: COLUMN exs_kline.start_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.start_time IS 'the kline start time';


--
-- Name: COLUMN exs_kline.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_kline.update_time IS 'the kline update time';


--
-- Name: exs_kline_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_kline_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_kline_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_kline_tid_seq OWNED BY exs_kline.tid;


--
-- Name: exs_order; Type: TABLE; Schema: public;
--

CREATE TABLE exs_order (
    tid bigint NOT NULL,
    order_id character varying(64) NOT NULL,
    type integer NOT NULL,
    user_id bigint NOT NULL,
    creator bigint NOT NULL,
    symbol character varying(16) NOT NULL,
    side character varying(8) NOT NULL,
    quantity double precision DEFAULT 0 NOT NULL,
    filled double precision DEFAULT 0 NOT NULL,
    price double precision DEFAULT 0 NOT NULL,
    trigger_type integer DEFAULT 0 NOT NULL,
    trigger_price double precision DEFAULT 0 NOT NULL,
    avg_price double precision DEFAULT 0 NOT NULL,
    total_price double precision DEFAULT 0 NOT NULL,
    holding double precision DEFAULT 0 NOT NULL,
    profit double precision DEFAULT 0 NOT NULL,
    owned double precision DEFAULT 0 NOT NULL,
    unhedged double precision DEFAULT 0 NOT NULL,
    in_balance character varying(30) NOT NULL,
    in_filled double precision NOT NULL,
    out_balance character varying(30) NOT NULL,
    out_filled double precision NOT NULL,
    fee_balance character varying(30) NOT NULL,
    fee_filled double precision DEFAULT 0 NOT NULL,
    fee_rate double precision DEFAULT 0 NOT NULL,
    transaction jsonb DEFAULT '{}'::jsonb NOT NULL,
    fee_settled_status integer DEFAULT 0 NOT NULL,
    fee_settled_next timestamp with time zone NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_order.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.tid IS 'the primary key';


--
-- Name: COLUMN exs_order.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.order_id IS 'the order string id';


--
-- Name: COLUMN exs_order.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.type IS 'the order type, Trade=100: is trade type, Trigger=200: is trigger trade order, Blowup=300: is blow up type';


--
-- Name: COLUMN exs_order.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.user_id IS 'the order user id';


--
-- Name: COLUMN exs_order.creator; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.creator IS 'the order creator user id';


--
-- Name: COLUMN exs_order.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.symbol IS 'the order symbol';


--
-- Name: COLUMN exs_order.side; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.side IS 'the order side, Buy=buy: is buy side, Sell=sell: is sell side';


--
-- Name: COLUMN exs_order.quantity; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.quantity IS 'the order expected quantity';


--
-- Name: COLUMN exs_order.filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.filled IS 'the order filled quantity';


--
-- Name: COLUMN exs_order.price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.price IS 'the order expected price';


--
-- Name: COLUMN exs_order.trigger_type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.trigger_type IS 'the order trigger type, None=0:is none type, StopProfit=100: is stop profit type, StopLoss=200: is stop loss';


--
-- Name: COLUMN exs_order.trigger_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.trigger_price IS 'the order trigger price';


--
-- Name: COLUMN exs_order.avg_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.avg_price IS 'the order filled avg price';


--
-- Name: COLUMN exs_order.total_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.total_price IS 'the order filled total price';


--
-- Name: COLUMN exs_order.holding; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.holding IS 'the order holding';


--
-- Name: COLUMN exs_order.profit; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.profit IS 'the order profit';


--
-- Name: COLUMN exs_order.owned; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.owned IS 'the order owned count';


--
-- Name: COLUMN exs_order.unhedged; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.unhedged IS 'the order owned is unbalanced';


--
-- Name: COLUMN exs_order.in_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.in_balance IS 'the in balance asset key';


--
-- Name: COLUMN exs_order.in_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.in_filled IS 'the in balance filled amount';


--
-- Name: COLUMN exs_order.out_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.out_balance IS 'the out balance asset key';


--
-- Name: COLUMN exs_order.out_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.out_filled IS 'the out balance filled amount';


--
-- Name: COLUMN exs_order.fee_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.fee_balance IS 'the fee balance asset key';


--
-- Name: COLUMN exs_order.fee_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.fee_filled IS 'the fee amount';


--
-- Name: COLUMN exs_order.fee_rate; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.fee_rate IS 'the order fee rate';


--
-- Name: COLUMN exs_order.transaction; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.transaction IS 'the order transaction info';


--
-- Name: COLUMN exs_order.fee_settled_status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.fee_settled_status IS 'the order transaction detail';


--
-- Name: COLUMN exs_order.fee_settled_next; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.fee_settled_next IS 'the fee settled time';


--
-- Name: COLUMN exs_order.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.update_time IS 'the order update time';


--
-- Name: COLUMN exs_order.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.create_time IS 'the order create time';


--
-- Name: COLUMN exs_order.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order.status IS 'the order status, Waiting=100, Pending=200:is pending, Partialled=300:is partialled, Done=400:is done, PartCanceled=410: is partialled canceled, Canceled=420: is canceled';


--
-- Name: exs_order_comm; Type: TABLE; Schema: public;
--

CREATE TABLE exs_order_comm (
    tid bigint NOT NULL,
    order_id bigint NOT NULL,
    user_id bigint NOT NULL,
    type integer NOT NULL,
    in_balance character varying(30) NOT NULL,
    in_fee double precision NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_order_comm.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.tid IS 'the primary key';


--
-- Name: COLUMN exs_order_comm.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.order_id IS 'the order id';


--
-- Name: COLUMN exs_order_comm.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.user_id IS 'the user id';


--
-- Name: COLUMN exs_order_comm.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.type IS 'the comm type, Normal=100:is normal type';


--
-- Name: COLUMN exs_order_comm.in_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.in_balance IS 'the in balance asset key';


--
-- Name: COLUMN exs_order_comm.in_fee; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.in_fee IS 'the in balance fee';


--
-- Name: COLUMN exs_order_comm.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.create_time IS 'the comm create time';


--
-- Name: COLUMN exs_order_comm.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_order_comm.status IS 'the comm status, Normal=100:is normal';


--
-- Name: exs_order_comm_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_order_comm_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_order_comm_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_order_comm_tid_seq OWNED BY exs_order_comm.tid;


--
-- Name: exs_order_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_order_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_order_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_order_tid_seq OWNED BY exs_order.tid;


--
-- Name: exs_user; Type: TABLE; Schema: public;
--

CREATE TABLE exs_user (
    tid bigint NOT NULL,
    type smallint DEFAULT 100 NOT NULL,
    role smallint DEFAULT 100 NOT NULL,
    name character varying(255),
    account character varying(255),
    phone character varying(255),
    password character varying(255),
    trade_pass character varying(255),
    image text,
    fee jsonb DEFAULT '{}'::jsonb NOT NULL,
    external jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status smallint DEFAULT 100 NOT NULL
);


--
-- Name: COLUMN exs_user.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.tid IS 'the primary key';


--
-- Name: COLUMN exs_user.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.type IS 'the user type,Admin=10:is admin user, Normal=100:is normal user';


--
-- Name: COLUMN exs_user.role; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.role IS 'ther user role, Normal=100:is normal, Staff=200:is staff';


--
-- Name: COLUMN exs_user.name; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.name IS 'the user name';


--
-- Name: COLUMN exs_user.account; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.account IS 'the user account to login';


--
-- Name: COLUMN exs_user.phone; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.phone IS 'the user phone number to login';


--
-- Name: COLUMN exs_user.password; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.password IS 'the user password to login';


--
-- Name: COLUMN exs_user.trade_pass; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.trade_pass IS 'the user trade password';


--
-- Name: COLUMN exs_user.image; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.image IS 'the user image';


--
-- Name: COLUMN exs_user.fee; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.fee IS 'the user fee';


--
-- Name: COLUMN exs_user.external; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.external IS 'the user external info';


--
-- Name: COLUMN exs_user.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.update_time IS 'the last updat time';


--
-- Name: COLUMN exs_user.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.create_time IS 'the craete time';


--
-- Name: COLUMN exs_user.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_user.status IS 'the user status, Normal=100:is normal, Locked=200:is locked, Removed=-1:is deleted';


--
-- Name: exs_user_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE exs_user_tid_seq
    START WITH 100000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exs_user_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE exs_user_tid_seq OWNED BY exs_user.tid;


--
-- Name: exs_withdraw; Type: TABLE; Schema: public;
--

CREATE TABLE exs_withdraw (
    tid bigint DEFAULT nextval('exs_order_tid_seq'::regclass) NOT NULL,
    order_id character varying(64) NOT NULL,
    type integer NOT NULL,
    user_id bigint NOT NULL,
    creator bigint NOT NULL,
    asset character varying(16) NOT NULL,
    quantity double precision DEFAULT 0 NOT NULL,
    transaction jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp(6) with time zone NOT NULL,
    create_time timestamp(6) with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN exs_withdraw.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.tid IS 'the primary key';


--
-- Name: COLUMN exs_withdraw.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.order_id IS 'the withdraw order string id';


--
-- Name: COLUMN exs_withdraw.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.type IS 'the withdraw order type, Withdraw=100: is withdraw type, Topup=200: is topup type, Goldbar=300: is goldbar bar';


--
-- Name: COLUMN exs_withdraw.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.user_id IS 'the withdraw order user id';


--
-- Name: COLUMN exs_withdraw.creator; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.creator IS 'the withdraw order creator user id';


--
-- Name: COLUMN exs_withdraw.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.asset IS 'the withdraw asset';


--
-- Name: COLUMN exs_withdraw.quantity; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.quantity IS 'the withdraw order quantity';


--
-- Name: COLUMN exs_withdraw.transaction; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.transaction IS 'the withdraw order transaction info';


--
-- Name: COLUMN exs_withdraw.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.update_time IS 'the withdraw order update time';


--
-- Name: COLUMN exs_withdraw.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.create_time IS 'the withdraw order create time';


--
-- Name: COLUMN exs_withdraw.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN exs_withdraw.status IS 'the withdraw order status, Pending=100:is pending, Confirmed=200:is confirmed, Done=300:is done, Canceled=320: is canceled';


--
-- Name: exs_balance tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_balance ALTER COLUMN tid SET DEFAULT nextval('exs_balance_tid_seq'::regclass);


--
-- Name: exs_balance_history tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_balance_history ALTER COLUMN tid SET DEFAULT nextval('exs_balance_record_tid_seq'::regclass);


--
-- Name: exs_holding tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_holding ALTER COLUMN tid SET DEFAULT nextval('exs_holding_tid_seq'::regclass);


--
-- Name: exs_kline tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_kline ALTER COLUMN tid SET DEFAULT nextval('exs_kline_tid_seq'::regclass);


--
-- Name: exs_order tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_order ALTER COLUMN tid SET DEFAULT nextval('exs_order_tid_seq'::regclass);


--
-- Name: exs_order_comm tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_order_comm ALTER COLUMN tid SET DEFAULT nextval('exs_order_comm_tid_seq'::regclass);


--
-- Name: exs_user tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_user ALTER COLUMN tid SET DEFAULT nextval('exs_user_tid_seq'::regclass);


--
-- Name: exs_balance exs_balance_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_balance
    ADD CONSTRAINT exs_balance_pkey PRIMARY KEY (tid);


--
-- Name: exs_balance_history exs_balance_record_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_balance_history
    ADD CONSTRAINT exs_balance_record_pkey PRIMARY KEY (tid);


--
-- Name: exs_holding exs_holding_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_holding
    ADD CONSTRAINT exs_holding_pkey PRIMARY KEY (tid);


--
-- Name: exs_kline exs_kline_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_kline
    ADD CONSTRAINT exs_kline_pkey PRIMARY KEY (tid);


--
-- Name: exs_order_comm exs_order_comm_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_order_comm
    ADD CONSTRAINT exs_order_comm_pkey PRIMARY KEY (tid);


--
-- Name: exs_withdraw exs_order_copy1_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_withdraw
    ADD CONSTRAINT exs_order_copy1_pkey PRIMARY KEY (tid);


--
-- Name: exs_order exs_order_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_order
    ADD CONSTRAINT exs_order_pkey PRIMARY KEY (tid);


--
-- Name: exs_user exs_user_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY exs_user
    ADD CONSTRAINT exs_user_pkey PRIMARY KEY (tid);


--
-- Name: exs_balance_history_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_balance_history_status_idx ON exs_balance_history USING btree (status);


--
-- Name: exs_balance_history_user_asset_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_balance_history_user_asset_idx ON exs_balance_history USING btree (user_id, asset, create_time);


--
-- Name: exs_balance_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_balance_status_idx ON exs_balance USING btree (status);


--
-- Name: exs_balance_user_area_asset_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_balance_user_area_asset_idx ON exs_balance USING btree (user_id, area, asset);


--
-- Name: exs_holding_amount_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_holding_amount_idx ON exs_holding USING btree (amount);


--
-- Name: exs_holding_blowup_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_holding_blowup_idx ON exs_holding USING btree (blowup);


--
-- Name: exs_holding_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_holding_status_idx ON exs_holding USING btree (status);


--
-- Name: exs_holding_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_holding_update_time_idx ON exs_holding USING btree (update_time);


--
-- Name: exs_holding_user_symbol_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_holding_user_symbol_idx ON exs_holding USING btree (user_id, symbol);


--
-- Name: exs_kline_interval_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_kline_interval_idx ON exs_kline USING btree (interv);


--
-- Name: exs_kline_start_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_kline_start_time_idx ON exs_kline USING btree (start_time);


--
-- Name: exs_kline_symbol_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_kline_symbol_idx ON exs_kline USING btree (symbol);


--
-- Name: exs_order_comm_create_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_comm_create_time_idx ON exs_order_comm USING btree (create_time);


--
-- Name: exs_order_comm_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_comm_status_idx ON exs_order_comm USING btree (status);


--
-- Name: exs_order_comm_user_type_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_order_comm_user_type_idx ON exs_order_comm USING btree (order_id, user_id, type);


--
-- Name: exs_order_fee_settled_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_fee_settled_idx ON exs_order USING btree (fee_settled_status, fee_settled_next);


--
-- Name: exs_order_order_id_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_order_order_id_idx ON exs_order USING btree (order_id);


--
-- Name: exs_order_side_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_side_idx ON exs_order USING btree (side);


--
-- Name: exs_order_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_status_idx ON exs_order USING btree (status);


--
-- Name: exs_order_symobl_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_symobl_idx ON exs_order USING btree (symbol);


--
-- Name: exs_order_trigger_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_trigger_idx ON exs_order USING btree (trigger_type, trigger_price);


--
-- Name: exs_order_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_type_idx ON exs_order USING btree (type);


--
-- Name: exs_order_unhedged_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_unhedged_idx ON exs_order USING btree (unhedged);


--
-- Name: exs_order_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_update_time_idx ON exs_order USING btree (update_time);


--
-- Name: exs_order_user_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_order_user_id_idx ON exs_order USING btree (user_id);


--
-- Name: exs_user_account_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_user_account_idx ON exs_user USING btree (account);


--
-- Name: exs_user_password_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_password_idx ON exs_user USING btree (password);


--
-- Name: exs_user_phone_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_phone_idx ON exs_user USING btree (phone);


--
-- Name: exs_user_role_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_role_idx ON exs_user USING btree (role);


--
-- Name: exs_user_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_status_idx ON exs_user USING btree (status);


--
-- Name: exs_user_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_type_idx ON exs_user USING btree (type);


--
-- Name: exs_user_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_user_update_time_idx ON exs_user USING btree (update_time);


--
-- Name: exs_withdraw_asset_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_withdraw_asset_idx ON exs_withdraw USING btree (asset);


--
-- Name: exs_withdraw_order_id_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX exs_withdraw_order_id_idx ON exs_withdraw USING btree (order_id);


--
-- Name: exs_withdraw_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_withdraw_status_idx ON exs_withdraw USING btree (status);


--
-- Name: exs_withdraw_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_withdraw_type_idx ON exs_withdraw USING btree (type);


--
-- Name: exs_withdraw_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_withdraw_update_time_idx ON exs_withdraw USING btree (update_time);


--
-- Name: exs_withdraw_user_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX exs_withdraw_user_id_idx ON exs_withdraw USING btree (user_id);


--
-- PostgreSQL database dump complete
--

