--
-- PostgreSQL database dump
--

-- Dumped from database version 13.8 (Debian 13.8-1.pgdg110+1)
-- Dumped by pg_dump version 13.8 (Debian 13.8-1.pgdg110+1)


DROP INDEX IF EXISTS gex_withdraw_user_id_idx;
DROP INDEX IF EXISTS gex_withdraw_update_time_idx;
DROP INDEX IF EXISTS gex_withdraw_type_idx;
DROP INDEX IF EXISTS gex_withdraw_status_idx;
DROP INDEX IF EXISTS gex_withdraw_processed_idx;
DROP INDEX IF EXISTS gex_withdraw_order_id_idx;
DROP INDEX IF EXISTS gex_withdraw_asset_idx;
DROP INDEX IF EXISTS gex_wallet_user_method_idx;
DROP INDEX IF EXISTS gex_wallet_status_idx;
DROP INDEX IF EXISTS gex_wallet_method_address_idx;
DROP INDEX IF EXISTS gex_user_update_time_idx;
DROP INDEX IF EXISTS gex_user_type_idx;
DROP INDEX IF EXISTS gex_user_status_idx;
DROP INDEX IF EXISTS gex_user_role_idx;
DROP INDEX IF EXISTS gex_user_recrod_status_idx;
DROP INDEX IF EXISTS gex_user_record_user_id_idx;
DROP INDEX IF EXISTS gex_user_record_update_time_idx;
DROP INDEX IF EXISTS gex_user_record_type_idx;
DROP INDEX IF EXISTS gex_user_phone_idx;
DROP INDEX IF EXISTS gex_user_password_idx;
DROP INDEX IF EXISTS gex_user_email_idx;
DROP INDEX IF EXISTS gex_user_account_idx;
DROP INDEX IF EXISTS gex_order_user_id_idx;
DROP INDEX IF EXISTS gex_order_update_time_idx;
DROP INDEX IF EXISTS gex_order_unhedged_idx;
DROP INDEX IF EXISTS gex_order_type_idx;
DROP INDEX IF EXISTS gex_order_trigger_time_idx;
DROP INDEX IF EXISTS gex_order_trigger_price_idx;
DROP INDEX IF EXISTS gex_order_symobl_idx;
DROP INDEX IF EXISTS gex_order_status_idx;
DROP INDEX IF EXISTS gex_order_side_idx;
DROP INDEX IF EXISTS gex_order_order_id_idx;
DROP INDEX IF EXISTS gex_order_fee_settled_idx;
DROP INDEX IF EXISTS gex_order_comm_user_type_idx;
DROP INDEX IF EXISTS gex_order_comm_status_idx;
DROP INDEX IF EXISTS gex_order_comm_create_time_idx;
DROP INDEX IF EXISTS gex_order_area_idx;
DROP INDEX IF EXISTS gex_kline_symbol_idx;
DROP INDEX IF EXISTS gex_kline_start_time_idx;
DROP INDEX IF EXISTS gex_kline_interval_idx;
DROP INDEX IF EXISTS gex_holding_user_symbol_idx;
DROP INDEX IF EXISTS gex_holding_update_time_idx;
DROP INDEX IF EXISTS gex_holding_status_idx;
DROP INDEX IF EXISTS gex_holding_blowup_idx;
DROP INDEX IF EXISTS gex_holding_amount_idx;
DROP INDEX IF EXISTS gex_balance_user_area_asset_idx;
DROP INDEX IF EXISTS gex_balance_status_idx;
DROP INDEX IF EXISTS gex_balance_record_update_time_idx;
DROP INDEX IF EXISTS gex_balance_record_type_idx;
DROP INDEX IF EXISTS gex_balance_record_balance_id_idx;
DROP INDEX IF EXISTS gex_balance_history_user_asset_idx;
DROP INDEX IF EXISTS gex_balance_history_status_idx;
ALTER TABLE IF EXISTS gex_wallet ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_user_record ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_user ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_order_comm ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_order ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_message ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_kline ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_holding ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_balance_record ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_balance_history ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS gex_balance ALTER COLUMN tid DROP DEFAULT;
DROP TABLE IF EXISTS gex_withdraw;
DROP SEQUENCE IF EXISTS gex_wallet_tid_seq;
DROP TABLE IF EXISTS gex_wallet;
DROP SEQUENCE IF EXISTS gex_user_tid_seq;
DROP SEQUENCE IF EXISTS gex_user_record_tid_seq;
DROP TABLE IF EXISTS gex_user_record;
DROP TABLE IF EXISTS gex_user;
DROP SEQUENCE IF EXISTS gex_order_tid_seq;
DROP SEQUENCE IF EXISTS gex_order_comm_tid_seq;
DROP TABLE IF EXISTS gex_order_comm;
DROP TABLE IF EXISTS gex_order;
DROP SEQUENCE IF EXISTS gex_message_tid_seq;
DROP TABLE IF EXISTS gex_message;
DROP SEQUENCE IF EXISTS gex_kline_tid_seq;
DROP TABLE IF EXISTS gex_kline;
DROP SEQUENCE IF EXISTS gex_holding_tid_seq;
DROP TABLE IF EXISTS gex_holding;
DROP SEQUENCE IF EXISTS gex_balance_tid_seq;
DROP SEQUENCE IF EXISTS gex_balance_record_tid_seq;
DROP TABLE IF EXISTS gex_balance_record;
DROP SEQUENCE IF EXISTS gex_balance_history_tid_seq;
DROP TABLE IF EXISTS gex_balance_history;
DROP TABLE IF EXISTS gex_balance;


--
-- Name: gex_balance; Type: TABLE; Schema: public;
--

CREATE TABLE gex_balance (
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
-- Name: COLUMN gex_balance.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.tid IS 'the primary key';


--
-- Name: COLUMN gex_balance.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.user_id IS 'the balance user id';


--
-- Name: COLUMN gex_balance.area; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.area IS 'the balance area, Funds=100:is funds area, Spot=200:is spot area, Futures=300:is futures area';


--
-- Name: COLUMN gex_balance.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.asset IS 'the balance asset key';


--
-- Name: COLUMN gex_balance.free; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.free IS 'the balance free amount';


--
-- Name: COLUMN gex_balance.locked; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.locked IS 'the balance locked amount';


--
-- Name: COLUMN gex_balance.margin; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.margin IS 'the balance margin value';


--
-- Name: COLUMN gex_balance.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.update_time IS 'the balance last update time';


--
-- Name: COLUMN gex_balance.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.create_time IS 'the balance create time';


--
-- Name: COLUMN gex_balance.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance.status IS 'the balance status, Normal=100: is normal, Locked=200: is locked';


--
-- Name: gex_balance_history; Type: TABLE; Schema: public;
--

CREATE TABLE gex_balance_history (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    asset character varying(30) NOT NULL,
    valuation double precision DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_balance_history.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.tid IS 'the primary key';


--
-- Name: COLUMN gex_balance_history.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.user_id IS 'the balance user id';


--
-- Name: COLUMN gex_balance_history.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.asset IS 'the balance asset key';


--
-- Name: COLUMN gex_balance_history.valuation; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.valuation IS 'the balance valuation';


--
-- Name: COLUMN gex_balance_history.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.update_time IS 'the balance record update time';


--
-- Name: COLUMN gex_balance_history.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.create_time IS 'the balance record create time, is daily zero time';


--
-- Name: COLUMN gex_balance_history.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_history.status IS 'the balance record status, Normal=100: is normal status';


--
-- Name: gex_balance_history_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_balance_history_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_balance_history_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_balance_history_tid_seq OWNED BY gex_balance_history.tid;


--
-- Name: gex_balance_record; Type: TABLE; Schema: public;
--

CREATE TABLE gex_balance_record (
    tid bigint NOT NULL,
    creator bigint NOT NULL,
    balance_id bigint NOT NULL,
    type integer NOT NULL,
    source character varying(64),
    target integer DEFAULT 0 NOT NULL,
    changed double precision NOT NULL,
    transaction jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp(6) with time zone NOT NULL,
    create_time timestamp(6) with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_balance_record.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.tid IS 'the primary key';


--
-- Name: COLUMN gex_balance_record.creator; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.creator IS 'the balance creator';


--
-- Name: COLUMN gex_balance_record.balance_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.balance_id IS 'the balance id';


--
-- Name: COLUMN gex_balance_record.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.type IS 'the balance record type, Trade=100: is trade type, TradeFee=110:is trade fee, Profit=200:is close profit, Blowup=210:is blowup, Transfer=300:is transfer, TransferInner=310:is transfer inner, Change=400: is manual change type, Topup=500: is topup, Withdraw=600: is withdraw, Goldbar=700:is gold bar';


--
-- Name: COLUMN gex_balance_record.source; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.source IS 'the balance record source';


--
-- Name: COLUMN gex_balance_record.target; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.target IS 'the balance target type';


--
-- Name: COLUMN gex_balance_record.changed; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.changed IS 'the balance change value';


--
-- Name: COLUMN gex_balance_record.transaction; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.transaction IS 'the balance record transaction info';


--
-- Name: COLUMN gex_balance_record.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.update_time IS 'the balance last update time';


--
-- Name: COLUMN gex_balance_record.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.create_time IS 'the balance create time';


--
-- Name: COLUMN gex_balance_record.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_balance_record.status IS 'the balance status, Pending=90:is pending, Normal=100: is normal, Fail=110: is fail.';


--
-- Name: gex_balance_record_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_balance_record_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_balance_record_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_balance_record_tid_seq OWNED BY gex_balance_record.tid;


--
-- Name: gex_balance_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_balance_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_balance_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_balance_tid_seq OWNED BY gex_balance.tid;


--
-- Name: gex_holding; Type: TABLE; Schema: public;
--

CREATE TABLE gex_holding (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    symbol character varying(16) NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    open double precision DEFAULT 0 NOT NULL,
    blowup double precision DEFAULT 0 NOT NULL,
    lever integer DEFAULT 5 NOT NULL,
    margin_used double precision DEFAULT 0 NOT NULL,
    margin_added double precision DEFAULT 0 NOT NULL,
    update_time timestamp(6) with time zone NOT NULL,
    create_time timestamp(6) with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_holding.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.tid IS 'the primary key';


--
-- Name: COLUMN gex_holding.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.user_id IS 'the holding user id';


--
-- Name: COLUMN gex_holding.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.symbol IS 'the holding symbol';


--
-- Name: COLUMN gex_holding.amount; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.amount IS 'the holding amount';


--
-- Name: COLUMN gex_holding.open; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.open IS 'the holding open price';


--
-- Name: COLUMN gex_holding.blowup; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.blowup IS 'the holding blowup price';


--
-- Name: COLUMN gex_holding.lever; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.lever IS 'the holding lever';


--
-- Name: COLUMN gex_holding.margin_used; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.margin_used IS 'the holding margin used';


--
-- Name: COLUMN gex_holding.margin_added; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.margin_added IS 'the holding margin added';


--
-- Name: COLUMN gex_holding.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.update_time IS 'the holding last update time';


--
-- Name: COLUMN gex_holding.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.create_time IS 'the holding create time';


--
-- Name: COLUMN gex_holding.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_holding.status IS 'the holding status, Normal=100: is normal, Locked=200: is locked';


--
-- Name: gex_holding_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_holding_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_holding_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_holding_tid_seq OWNED BY gex_holding.tid;


--
-- Name: gex_kline; Type: TABLE; Schema: public;
--

CREATE TABLE gex_kline (
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
-- Name: COLUMN gex_kline.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.tid IS 'the primay key';


--
-- Name: COLUMN gex_kline.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.symbol IS 'the kline symbol';


--
-- Name: COLUMN gex_kline.interv; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.interv IS 'the kline interval key';


--
-- Name: COLUMN gex_kline.amount; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.amount IS 'the kline amount';


--
-- Name: COLUMN gex_kline.count; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.count IS 'the kline count';


--
-- Name: COLUMN gex_kline.open; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.open IS 'the kline open price';


--
-- Name: COLUMN gex_kline.close; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.close IS 'the kline close price';


--
-- Name: COLUMN gex_kline.low; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.low IS 'the kline low price';


--
-- Name: COLUMN gex_kline.high; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.high IS 'the kline high price';


--
-- Name: COLUMN gex_kline.volume; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.volume IS 'the kline volume price';


--
-- Name: COLUMN gex_kline.start_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.start_time IS 'the kline start time';


--
-- Name: COLUMN gex_kline.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_kline.update_time IS 'the kline update time';


--
-- Name: gex_kline_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_kline_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_kline_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_kline_tid_seq OWNED BY gex_kline.tid;


--
-- Name: gex_message; Type: TABLE; Schema: public;
--

CREATE TABLE gex_message (
    tid bigint NOT NULL,
    type integer NOT NULL,
    title jsonb DEFAULT '{}'::jsonb NOT NULL,
    content jsonb NOT NULL,
    to_user_id bigint NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_message.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.tid IS 'the primary key';


--
-- Name: COLUMN gex_message.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.type IS 'the message type, User=100:is user type, Global=200:is global type';


--
-- Name: COLUMN gex_message.title; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.title IS 'the message title';


--
-- Name: COLUMN gex_message.content; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.content IS 'the message content';


--
-- Name: COLUMN gex_message.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.update_time IS 'the message update time';


--
-- Name: COLUMN gex_message.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.create_time IS 'the message create time';


--
-- Name: COLUMN gex_message.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_message.status IS 'the message status, Normal=100: is normal status, Removed=-1:is removed';


--
-- Name: gex_message_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_message_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_message_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_message_tid_seq OWNED BY gex_message.tid;


--
-- Name: gex_order; Type: TABLE; Schema: public;
--

CREATE TABLE gex_order (
    tid bigint NOT NULL,
    order_id character varying(64) NOT NULL,
    type integer NOT NULL,
    user_id bigint NOT NULL,
    creator bigint NOT NULL,
    area integer DEFAULT 0 NOT NULL,
    symbol character varying(16) NOT NULL,
    side character varying(8) NOT NULL,
    quantity double precision DEFAULT 0 NOT NULL,
    filled double precision DEFAULT 0 NOT NULL,
    price double precision DEFAULT 0 NOT NULL,
    trigger_type integer DEFAULT 0 NOT NULL,
    trigger_price double precision DEFAULT 0 NOT NULL,
    trigger_time timestamp with time zone NOT NULL,
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
-- Name: COLUMN gex_order.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.tid IS 'the primary key';


--
-- Name: COLUMN gex_order.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.order_id IS 'the order string id';


--
-- Name: COLUMN gex_order.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.type IS 'the order type, Trade=100: is trade type, Trigger=200: is trigger trade order, Blowup=300: is blow up type';


--
-- Name: COLUMN gex_order.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.user_id IS 'the order user id';


--
-- Name: COLUMN gex_order.creator; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.creator IS 'the order creator user id';


--
-- Name: COLUMN gex_order.area; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.area IS 'the order area, None=0, Spot=OrderArea(BalanceAreaSpot):is spot area, Futures=OrderArea(BalanceAreaFutures):is futures area';


--
-- Name: COLUMN gex_order.symbol; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.symbol IS 'the order symbol';


--
-- Name: COLUMN gex_order.side; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.side IS 'the order side, Buy=buy: is buy side, Sell=sell: is sell side';


--
-- Name: COLUMN gex_order.quantity; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.quantity IS 'the order expected quantity';


--
-- Name: COLUMN gex_order.filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.filled IS 'the order filled quantity';


--
-- Name: COLUMN gex_order.price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.price IS 'the order expected price';


--
-- Name: COLUMN gex_order.trigger_type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.trigger_type IS 'the order trigger type, None=0:is none type, StopProfit=100: is stop profit type, StopLoss=200: is stop loss, AfterOpen=300: is after open, AfterClose=310: is after close';


--
-- Name: COLUMN gex_order.trigger_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.trigger_price IS 'the order trigger price';


--
-- Name: COLUMN gex_order.trigger_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.trigger_time IS 'the order trigger time';


--
-- Name: COLUMN gex_order.avg_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.avg_price IS 'the order filled avg price';


--
-- Name: COLUMN gex_order.total_price; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.total_price IS 'the order filled total price';


--
-- Name: COLUMN gex_order.holding; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.holding IS 'the order holding';


--
-- Name: COLUMN gex_order.profit; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.profit IS 'the order profit';


--
-- Name: COLUMN gex_order.owned; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.owned IS 'the order owned count';


--
-- Name: COLUMN gex_order.unhedged; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.unhedged IS 'the order owned is unbalanced';


--
-- Name: COLUMN gex_order.in_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.in_balance IS 'the in balance asset key';


--
-- Name: COLUMN gex_order.in_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.in_filled IS 'the in balance filled amount';


--
-- Name: COLUMN gex_order.out_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.out_balance IS 'the out balance asset key';


--
-- Name: COLUMN gex_order.out_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.out_filled IS 'the out balance filled amount';


--
-- Name: COLUMN gex_order.fee_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.fee_balance IS 'the fee balance asset key';


--
-- Name: COLUMN gex_order.fee_filled; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.fee_filled IS 'the fee amount';


--
-- Name: COLUMN gex_order.fee_rate; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.fee_rate IS 'the order fee rate';


--
-- Name: COLUMN gex_order.transaction; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.transaction IS 'the order transaction info';


--
-- Name: COLUMN gex_order.fee_settled_status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.fee_settled_status IS 'the order transaction detail';


--
-- Name: COLUMN gex_order.fee_settled_next; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.fee_settled_next IS 'the fee settled time';


--
-- Name: COLUMN gex_order.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.update_time IS 'the order update time';


--
-- Name: COLUMN gex_order.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.create_time IS 'the order create time';


--
-- Name: COLUMN gex_order.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order.status IS 'the order status, Waiting=100, Pending=200:is pending, Partialled=300:is partialled, Done=400:is done, PartCanceled=410: is partialled canceled, Canceled=420: is canceled';


--
-- Name: gex_order_comm; Type: TABLE; Schema: public;
--

CREATE TABLE gex_order_comm (
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
-- Name: COLUMN gex_order_comm.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.tid IS 'the primary key';


--
-- Name: COLUMN gex_order_comm.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.order_id IS 'the order id';


--
-- Name: COLUMN gex_order_comm.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.user_id IS 'the user id';


--
-- Name: COLUMN gex_order_comm.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.type IS 'the comm type, Normal=100:is normal type';


--
-- Name: COLUMN gex_order_comm.in_balance; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.in_balance IS 'the in balance asset key';


--
-- Name: COLUMN gex_order_comm.in_fee; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.in_fee IS 'the in balance fee';


--
-- Name: COLUMN gex_order_comm.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.create_time IS 'the comm create time';


--
-- Name: COLUMN gex_order_comm.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_order_comm.status IS 'the comm status, Normal=100:is normal';


--
-- Name: gex_order_comm_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_order_comm_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_order_comm_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_order_comm_tid_seq OWNED BY gex_order_comm.tid;


--
-- Name: gex_order_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_order_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_order_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_order_tid_seq OWNED BY gex_order.tid;


--
-- Name: gex_user; Type: TABLE; Schema: public;
--

CREATE TABLE gex_user (
    tid bigint NOT NULL,
    type smallint DEFAULT 100 NOT NULL,
    role smallint DEFAULT 100 NOT NULL,
    name character varying(255),
    account character varying(255),
    phone character varying(255),
    email character varying(255),
    password character varying(255),
    trade_pass character varying(255),
    image text,
    fee jsonb DEFAULT '{}'::jsonb NOT NULL,
    external jsonb DEFAULT '{}'::jsonb NOT NULL,
    favorites jsonb DEFAULT '{}'::jsonb NOT NULL,
    config jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status smallint DEFAULT 100 NOT NULL
);


--
-- Name: COLUMN gex_user.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.tid IS 'the primary key';


--
-- Name: COLUMN gex_user.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.type IS 'the user type,Admin=10:is admin user, Normal=100:is normal user';


--
-- Name: COLUMN gex_user.role; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.role IS 'ther user role, Normal=100:is normal, Staff=200:is staff, Maker=300:is maker';


--
-- Name: COLUMN gex_user.name; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.name IS 'the user name';


--
-- Name: COLUMN gex_user.account; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.account IS 'the user account to login';


--
-- Name: COLUMN gex_user.phone; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.phone IS 'the user phone number to login';


--
-- Name: COLUMN gex_user.email; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.email IS 'the user email';


--
-- Name: COLUMN gex_user.password; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.password IS 'the user password to login';


--
-- Name: COLUMN gex_user.trade_pass; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.trade_pass IS 'the user trade password';


--
-- Name: COLUMN gex_user.image; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.image IS 'the user image';


--
-- Name: COLUMN gex_user.fee; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.fee IS 'the user fee';


--
-- Name: COLUMN gex_user.external; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.external IS 'the user external info';


--
-- Name: COLUMN gex_user.favorites; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.favorites IS 'the user favorites';


--
-- Name: COLUMN gex_user.config; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.config IS 'the user config';


--
-- Name: COLUMN gex_user.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.update_time IS 'the last updat time';


--
-- Name: COLUMN gex_user.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.create_time IS 'the craete time';


--
-- Name: COLUMN gex_user.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user.status IS 'the user status, Normal=100:is normal, Locked=200:is locked, Removed=-1:is deleted';


--
-- Name: gex_user_record; Type: TABLE; Schema: public;
--

CREATE TABLE gex_user_record (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    type integer NOT NULL,
    from_addr character varying(255) NOT NULL,
    external jsonb DEFAULT '{}'::jsonb NOT NULL,
    prev_id bigint DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_user_record.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.tid IS 'the primary key';


--
-- Name: COLUMN gex_user_record.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.user_id IS 'the user record user id';


--
-- Name: COLUMN gex_user_record.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.type IS 'the user recrod type, Login=100:is login record type';


--
-- Name: COLUMN gex_user_record.from_addr; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.from_addr IS 'the user record from addr';


--
-- Name: COLUMN gex_user_record.external; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.external IS 'the user record external info';


--
-- Name: COLUMN gex_user_record.prev_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.prev_id IS 'the user record prev id';


--
-- Name: COLUMN gex_user_record.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.update_time IS 'the user recrod update time';


--
-- Name: COLUMN gex_user_record.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.create_time IS 'the user record create time';


--
-- Name: COLUMN gex_user_record.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_user_record.status IS 'the user record status, Normal=100:is normal status';


--
-- Name: gex_user_record_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_user_record_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_user_record_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_user_record_tid_seq OWNED BY gex_user_record.tid;


--
-- Name: gex_user_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_user_tid_seq
    START WITH 100000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_user_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_user_tid_seq OWNED BY gex_user.tid;


--
-- Name: gex_wallet; Type: TABLE; Schema: public;
--

CREATE TABLE gex_wallet (
    tid bigint NOT NULL,
    user_id bigint NOT NULL,
    method character varying(255) NOT NULL,
    address character varying(255) NOT NULL,
    update_time timestamp with time zone NOT NULL,
    create_time timestamp with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_wallet.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.tid IS 'the wallet primary key';


--
-- Name: COLUMN gex_wallet.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.user_id IS 'the wallet user id';


--
-- Name: COLUMN gex_wallet.method; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.method IS 'the wallet type, Tron=tron: is tron method, Ethereum=ethereum: is ethereum method';


--
-- Name: COLUMN gex_wallet.address; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.address IS 'the wallet address';


--
-- Name: COLUMN gex_wallet.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.update_time IS 'the wallet update time';


--
-- Name: COLUMN gex_wallet.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.create_time IS 'the wallet create time';


--
-- Name: COLUMN gex_wallet.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_wallet.status IS 'the wallet status, Normal=100:is normal';


--
-- Name: gex_wallet_tid_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE gex_wallet_tid_seq
    START WITH 1000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gex_wallet_tid_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE gex_wallet_tid_seq OWNED BY gex_wallet.tid;


--
-- Name: gex_withdraw; Type: TABLE; Schema: public;
--

CREATE TABLE gex_withdraw (
    tid bigint DEFAULT nextval('gex_order_tid_seq'::regclass) NOT NULL,
    order_id character varying(255) NOT NULL,
    type integer NOT NULL,
    user_id bigint NOT NULL,
    creator bigint NOT NULL,
    method character varying(255) NOT NULL,
    asset character varying(16) NOT NULL,
    quantity double precision DEFAULT 0 NOT NULL,
    sender character varying(255),
    receiver character varying(255) NOT NULL,
    processed integer DEFAULT 0 NOT NULL,
    result jsonb DEFAULT '{}'::jsonb NOT NULL,
    update_time timestamp(6) with time zone NOT NULL,
    create_time timestamp(6) with time zone NOT NULL,
    status integer NOT NULL
);


--
-- Name: COLUMN gex_withdraw.tid; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.tid IS 'the primary key';


--
-- Name: COLUMN gex_withdraw.order_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.order_id IS 'the withdraw order string id';


--
-- Name: COLUMN gex_withdraw.type; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.type IS 'the withdraw order type, Withdraw=100: is withdraw type, Topup=200: is topup type, Goldbar=300: is goldbar bar';


--
-- Name: COLUMN gex_withdraw.user_id; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.user_id IS 'the withdraw order user id';


--
-- Name: COLUMN gex_withdraw.creator; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.creator IS 'the withdraw order creator user id';


--
-- Name: COLUMN gex_withdraw.method; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.method IS 'the withdraw metod, Tron=tron: is tron method, Ethereum=ethereum: is ethereum method';


--
-- Name: COLUMN gex_withdraw.asset; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.asset IS 'the withdraw asset';


--
-- Name: COLUMN gex_withdraw.quantity; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.quantity IS 'the withdraw order quantity';


--
-- Name: COLUMN gex_withdraw.receiver; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.receiver IS 'the widhdraw receiver';


--
-- Name: COLUMN gex_withdraw.processed; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.processed IS 'the withdraw if processed';


--
-- Name: COLUMN gex_withdraw.result; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.result IS 'the withdraw order transaction info';


--
-- Name: COLUMN gex_withdraw.update_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.update_time IS 'the withdraw order update time';


--
-- Name: COLUMN gex_withdraw.create_time; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.create_time IS 'the withdraw order create time';


--
-- Name: COLUMN gex_withdraw.status; Type: COMMENT; Schema: public;
--

COMMENT ON COLUMN gex_withdraw.status IS 'the withdraw order status, Pending=100:is pending, Confirmed=200:is confirmed, Done=300:is done, Canceled=320: is canceled';


--
-- Name: gex_balance tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance ALTER COLUMN tid SET DEFAULT nextval('gex_balance_tid_seq'::regclass);


--
-- Name: gex_balance_history tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance_history ALTER COLUMN tid SET DEFAULT nextval('gex_balance_history_tid_seq'::regclass);


--
-- Name: gex_balance_record tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance_record ALTER COLUMN tid SET DEFAULT nextval('gex_balance_record_tid_seq'::regclass);


--
-- Name: gex_holding tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_holding ALTER COLUMN tid SET DEFAULT nextval('gex_holding_tid_seq'::regclass);


--
-- Name: gex_kline tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_kline ALTER COLUMN tid SET DEFAULT nextval('gex_kline_tid_seq'::regclass);


--
-- Name: gex_message tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_message ALTER COLUMN tid SET DEFAULT nextval('gex_message_tid_seq'::regclass);


--
-- Name: gex_order tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_order ALTER COLUMN tid SET DEFAULT nextval('gex_order_tid_seq'::regclass);


--
-- Name: gex_order_comm tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_order_comm ALTER COLUMN tid SET DEFAULT nextval('gex_order_comm_tid_seq'::regclass);


--
-- Name: gex_user tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_user ALTER COLUMN tid SET DEFAULT nextval('gex_user_tid_seq'::regclass);


--
-- Name: gex_user_record tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_user_record ALTER COLUMN tid SET DEFAULT nextval('gex_user_record_tid_seq'::regclass);


--
-- Name: gex_wallet tid; Type: DEFAULT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_wallet ALTER COLUMN tid SET DEFAULT nextval('gex_wallet_tid_seq'::regclass);


--
-- Name: gex_balance_history gex_balance_history_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance_history
    ADD CONSTRAINT gex_balance_history_pkey PRIMARY KEY (tid);


--
-- Name: gex_balance gex_balance_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance
    ADD CONSTRAINT gex_balance_pkey PRIMARY KEY (tid);


--
-- Name: gex_balance_record gex_balance_record_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_balance_record
    ADD CONSTRAINT gex_balance_record_pkey PRIMARY KEY (tid);


--
-- Name: gex_holding gex_holding_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_holding
    ADD CONSTRAINT gex_holding_pkey PRIMARY KEY (tid);


--
-- Name: gex_kline gex_kline_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_kline
    ADD CONSTRAINT gex_kline_pkey PRIMARY KEY (tid);


--
-- Name: gex_message gex_message_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_message
    ADD CONSTRAINT gex_message_pkey PRIMARY KEY (tid);


--
-- Name: gex_order_comm gex_order_comm_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_order_comm
    ADD CONSTRAINT gex_order_comm_pkey PRIMARY KEY (tid);


--
-- Name: gex_withdraw gex_order_copy1_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_withdraw
    ADD CONSTRAINT gex_order_copy1_pkey PRIMARY KEY (tid);


--
-- Name: gex_order gex_order_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_order
    ADD CONSTRAINT gex_order_pkey PRIMARY KEY (tid);


--
-- Name: gex_user gex_user_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_user
    ADD CONSTRAINT gex_user_pkey PRIMARY KEY (tid);


--
-- Name: gex_user_record gex_user_record_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_user_record
    ADD CONSTRAINT gex_user_record_pkey PRIMARY KEY (tid);


--
-- Name: gex_wallet gex_wallet_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE IF EXISTS ONLY gex_wallet
    ADD CONSTRAINT gex_wallet_pkey PRIMARY KEY (tid);


--
-- Name: gex_balance_history_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_balance_history_status_idx ON gex_balance_history USING btree (status);


--
-- Name: gex_balance_history_user_asset_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_balance_history_user_asset_idx ON gex_balance_history USING btree (user_id, asset, create_time);


--
-- Name: gex_balance_record_balance_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_balance_record_balance_id_idx ON gex_balance_record USING btree (balance_id);


--
-- Name: gex_balance_record_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_balance_record_type_idx ON gex_balance_record USING btree (type);


--
-- Name: gex_balance_record_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_balance_record_update_time_idx ON gex_balance_record USING btree (update_time);


--
-- Name: gex_balance_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_balance_status_idx ON gex_balance USING btree (status);


--
-- Name: gex_balance_user_area_asset_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_balance_user_area_asset_idx ON gex_balance USING btree (user_id, area, asset);


--
-- Name: gex_holding_amount_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_holding_amount_idx ON gex_holding USING btree (amount);


--
-- Name: gex_holding_blowup_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_holding_blowup_idx ON gex_holding USING btree (blowup);


--
-- Name: gex_holding_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_holding_status_idx ON gex_holding USING btree (status);


--
-- Name: gex_holding_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_holding_update_time_idx ON gex_holding USING btree (update_time);


--
-- Name: gex_holding_user_symbol_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_holding_user_symbol_idx ON gex_holding USING btree (user_id, symbol);


--
-- Name: gex_kline_interval_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_kline_interval_idx ON gex_kline USING btree (interv);


--
-- Name: gex_kline_start_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_kline_start_time_idx ON gex_kline USING btree (start_time);


--
-- Name: gex_kline_symbol_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_kline_symbol_idx ON gex_kline USING btree (symbol);


--
-- Name: gex_order_area_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_area_idx ON gex_order USING btree (area);


--
-- Name: gex_order_comm_create_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_comm_create_time_idx ON gex_order_comm USING btree (create_time);


--
-- Name: gex_order_comm_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_comm_status_idx ON gex_order_comm USING btree (status);


--
-- Name: gex_order_comm_user_type_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_order_comm_user_type_idx ON gex_order_comm USING btree (order_id, user_id, type);


--
-- Name: gex_order_fee_settled_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_fee_settled_idx ON gex_order USING btree (fee_settled_status, fee_settled_next);


--
-- Name: gex_order_order_id_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_order_order_id_idx ON gex_order USING btree (order_id);


--
-- Name: gex_order_side_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_side_idx ON gex_order USING btree (side);


--
-- Name: gex_order_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_status_idx ON gex_order USING btree (status);


--
-- Name: gex_order_symobl_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_symobl_idx ON gex_order USING btree (symbol);


--
-- Name: gex_order_trigger_price_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_trigger_price_idx ON gex_order USING btree (trigger_type, trigger_price);


--
-- Name: gex_order_trigger_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_trigger_time_idx ON gex_order USING btree (trigger_type, trigger_time);


--
-- Name: gex_order_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_type_idx ON gex_order USING btree (type);


--
-- Name: gex_order_unhedged_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_unhedged_idx ON gex_order USING btree (unhedged);


--
-- Name: gex_order_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_update_time_idx ON gex_order USING btree (update_time);


--
-- Name: gex_order_user_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_order_user_id_idx ON gex_order USING btree (user_id);


--
-- Name: gex_user_account_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_user_account_idx ON gex_user USING btree (account);


--
-- Name: gex_user_email_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_user_email_idx ON gex_user USING btree (email);


--
-- Name: gex_user_password_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_password_idx ON gex_user USING btree (password);


--
-- Name: gex_user_phone_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_user_phone_idx ON gex_user USING btree (phone);


--
-- Name: gex_user_record_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_record_type_idx ON gex_user_record USING btree (type);


--
-- Name: gex_user_record_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_record_update_time_idx ON gex_user_record USING btree (update_time);


--
-- Name: gex_user_record_user_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_record_user_id_idx ON gex_user_record USING btree (user_id);


--
-- Name: gex_user_recrod_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_recrod_status_idx ON gex_user_record USING btree (status);


--
-- Name: gex_user_role_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_role_idx ON gex_user USING btree (role);


--
-- Name: gex_user_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_status_idx ON gex_user USING btree (status);


--
-- Name: gex_user_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_type_idx ON gex_user USING btree (type);


--
-- Name: gex_user_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_user_update_time_idx ON gex_user USING btree (update_time);


--
-- Name: gex_wallet_method_address_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_wallet_method_address_idx ON gex_wallet USING btree (method, address);


--
-- Name: gex_wallet_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_wallet_status_idx ON gex_wallet USING btree (status);


--
-- Name: gex_wallet_user_method_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_wallet_user_method_idx ON gex_wallet USING btree (user_id, method);


--
-- Name: gex_withdraw_asset_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_asset_idx ON gex_withdraw USING btree (asset);


--
-- Name: gex_withdraw_order_id_idx; Type: INDEX; Schema: public;
--

CREATE UNIQUE INDEX gex_withdraw_order_id_idx ON gex_withdraw USING btree (order_id);


--
-- Name: gex_withdraw_processed_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_processed_idx ON gex_withdraw USING btree (processed);


--
-- Name: gex_withdraw_status_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_status_idx ON gex_withdraw USING btree (status);


--
-- Name: gex_withdraw_type_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_type_idx ON gex_withdraw USING btree (type);


--
-- Name: gex_withdraw_update_time_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_update_time_idx ON gex_withdraw USING btree (update_time);


--
-- Name: gex_withdraw_user_id_idx; Type: INDEX; Schema: public;
--

CREATE INDEX gex_withdraw_user_id_idx ON gex_withdraw USING btree (user_id);


--
-- PostgreSQL database dump complete
--

