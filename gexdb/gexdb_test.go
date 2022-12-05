package gexdb

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/Centny/rediscache"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xprop"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/gexupgrade"
)

var proxyAddr = "127.0.0.1:1105"
var ctx = context.Background()

func init() {
	func() {
		defer func() {
			recover()
		}()
		Pool()
	}()
	func() {
		defer func() {
			recover()
		}()
		Redis()
	}()
	rediscache.InitRedisPool("redis.loc:6379?db=11")
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/gexservice")
	if err != nil {
		panic(err)
	}
	Pool = pgx.Pool
	Redis = rediscache.C
	basedb.SYS = "gex"
	basedb.Pool = pgx.Pool
	_, _, err = Pool().Exec(ctx, gexupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, _, err = Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "gex_"))
	if err != nil {
		panic(err)
	}
	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	_, err = CheckDb(ctx)
	if err != nil {
		panic(err)
	}
	proxyServer := os.Getenv("PROXY_SERVER")
	if len(proxyServer) > 0 {
		proxyAddr = proxyServer
	}
	basedb.StoreConf(ctx, ConfigBrokerCommRate, "0.8")
	config := xprop.NewConfig()
	config.LoadPropString(`
[message.withdraw.done.title]
en=Withdraw Success
_=提现成功

[message.withdraw.done.content]
en=you withdraw ${_amount}${_asset} success on ${_time}
_=您于${_time}提现${_amount}${_asset}成功

[message.withdraw.fail.title]
en=Withdraw Fail
_=提现失败

[message.withdraw.fail.content]
en=you withdraw ${_amount}${_asset} Fail on ${_time}, ${_message}
_=您于${_time}提现${_amount}${_asset}失败

[message.topup.title]
en=Topup Success
_=充值成功

[message.topup.content]
en=you topup ${_amount}${_asset} success on ${_time}
_=您于${_time}成功充值${_amount}${_asset}

[message.blowup.title]
en=Blowup Warning
_=爆仓提醒

[message.blowup.content]
en=you position ${_amount} ${_symbol} is blowup on ${_time}, open price is ${_openPrice}, mark price is ${_markPrice}
_=您的仓位${_amount} ${_symbol}于${_time}爆仓，开仓价格为${_openPrice}，标记价格为${_markPrice}
	`)
	MessageTemplate = ReadMessageTemplateByConfig(config)
	go http.ListenAndServe(":6062", nil)
}

func clear() {
	_, _, err := Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.CLEAR, "_sys_", "gex_"))
	if err != nil {
		panic(err)
	}
	_, _, err = Pool().Exec(ctx, gexupgrade.CLEAR)
	if err != nil {
		panic(err)
	}
}
