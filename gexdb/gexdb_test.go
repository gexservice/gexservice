package gexdb

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/Centny/rediscache"
	"github.com/codingeasygo/crud/pgx"
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
	basedb.SYS = "exs"
	basedb.Pool = pgx.Pool
	_, _, err = Pool().Exec(ctx, gexupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, _, err = Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "exs_"))
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
	go http.ListenAndServe(":6062", nil)
}

func clear() {
	_, _, err := Pool().Exec(ctx, gexupgrade.CLEAR)
	if err != nil {
		panic(err)
	}
}
