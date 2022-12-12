package gexupgrade

const INIT = `
INSERT INTO gex_user(tid,type,name,account,password,update_time,create_time,status) VALUES (1000, 10, 'admin', 'admin', '40bd001563085fc35165329ea1ff5c5ecbdbbeef', '2021-06-15 09:34:37.664393+00', '2021-06-15 09:34:37.664393+00', 100);
INSERT INTO gex_config(key,value,update_time) VALUES ('goldbar_address', '[]', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('goldbar_explain', '', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('goldbar_rate', '1600', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('goldbar_fee', '0.005', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('goldbar_tips', '', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('welcome_message', 'welcom', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('withdraw_max', '50000', '2021-07-04 12:51:17.074424+00');
INSERT INTO gex_config(key,value,update_time) VALUES ('trade_rule', 'rule', '2021-07-04 12:51:17.074424+00');
`

const CHECK = `
ALTER TABLE gex_user ADD COLUMN IF NOT EXISTS config jsonb DEFAULT '{}'::jsonb NOT NULL;
ALTER TABLE gex_user ADD COLUMN IF NOT EXISTS email character varying(255);
CREATE UNIQUE INDEX IF NOT EXISTS gex_user_email_idx ON gex_user USING btree (email);
ALTER TABLE gex_balance_record ADD COLUMN IF NOT EXISTS transaction jsonb DEFAULT '{}'::jsonb NOT NULL;
ALTER TABLE gex_balance_record ADD COLUMN IF NOT EXISTS source character varying(64);
ALTER TABLE gex_order ADD COLUMN IF NOT EXISTS area integer DEFAULT 0 NOT NULL;
CREATE INDEX IF NOT EXISTS gex_order_area_idx ON gex_order USING btree (area);
`
