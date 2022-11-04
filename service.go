package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/Centny/rediscache"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/filter"
	"github.com/gexservice/gexservice/base/baseapi"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/session"
	"github.com/gexservice/gexservice/base/transport"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexapi"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Printf(`service %v version\n`, Version)
		return
	}
	confPath := "conf/gexservice.properties"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}
	fmt.Printf("Environ:\n")
	for _, env := range os.Environ() {
		fmt.Printf("  %v\n", env)
	}
	var err error
	conf := xprop.NewConfig()
	conf.Load(confPath)
	conf.Print()
	rediscache.InitRedisPool(conf.Str("/server/redis_con"))
	_, err = pgx.Bootstrap(conf.Str("/server/pg_con"))
	if err != nil {
		panic(err)
	}
	basedb.SYS = "gex"
	basedb.Pool = pgx.Pool
	gexdb.Pool = pgx.Pool
	gexdb.Redis = rediscache.C
	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	_, err = gexdb.CheckDb(context.Background())
	if err != nil {
		panic(err)
	}
	// if len(os.Args) > 2 {
	// 	switch os.Args[2] {
	// 	case "testdata":
	// 		xlog.Warnf("generate test data is starting")
	// 		err = GenerateTestData()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		xlog.Warnf("generate test data is done")
	// 		return
	// 	}
	// }
	gob.Register(xmap.M{})
	sb := session.NewDbSessionBuilder()
	sb.Redis = rediscache.C
	sb.ShowLog = false
	web.Shared.Builder = sb
	if pgTransport := conf.StrDef("", "/server/pg_transport"); len(pgTransport) > 0 {
		xlog.Warnf("pg transport is starting by %v /transport/pg", pgTransport)
		forward, err := transport.NewTransportH(pgTransport)
		if err != nil {
			xlog.Errorf("pg transport start by %v fail with %v", pgTransport, err)
		} else {
			web.Shared.Handle("^/transport/pg(\\?.*)?$", forward)
		}
	}
	if redisTransport := conf.StrDef("", "/server/redis_transport"); len(redisTransport) > 0 {
		xlog.Warnf("redis transport is starting by %v on /transport/redis", redisTransport)
		forward, err := transport.NewTransportH(redisTransport)
		if err != nil {
			xlog.Errorf("redis transport start by %v fail with %v", redisTransport, err)
		} else {
			web.Shared.Handle("^/transport/redis(\\?.*)?$", forward)
		}
	}
	//base handler
	web.Shared.Filter("^.*$", filter.NewAllCORS())
	web.Shared.FilterFunc("^/(index.html)?(\\?.*)?$", filter.NoCacheF)
	web.Shared.FilterFunc("^/(usr|pub)/.*$", filter.NoCacheF)
	web.Shared.StartMonitor()
	web.HandleFunc("^/adm/status(\\?.*)?$", func(hs *web.Session) web.Result {
		res := xmap.M{}
		res["http"], _ = web.Shared.State()
		return hs.SendJSON(res)
	})
	{ //mp config
		conf.Range("admin", func(key string, val interface{}) { gexapi.ConfAdminH[key] = val })
	}
	// if conf.StrDef("", "mock/enabled") == "1" { //mock
	// 	xlog.Warnf("mock is enabled")
	// 	web.Shared.HandleFunc("^/mock/payTopupOrder(\\?.*)?$", gexapi.MockPayTopupOrderH)
	// }
	//
	err = matcher.Bootstrap(conf)
	if err != nil {
		panic(err)
	}
	market.Bootstrap()
	gexapi.Handle("", web.Shared)
	uploader := baseapi.NewUploadH(conf.StrDef("upload", "/server/upload"), "/upload")
	web.Handle("^/usr/upload(\\?.*)?$", uploader)
	web.Shared.HandleNormal("^/upload.*$", http.StripPrefix("/upload", http.FileServer(http.Dir(conf.StrDef("upload", "/server/upload")))))
	web.Shared.HandleNormal("^/debug/.*$", http.DefaultServeMux)
	web.Shared.HandleNormal("^/apidoc.*$", http.FileServer(http.Dir(conf.StrDef("www", "/www/apidoc"))))
	wwwFS := http.FileServer(http.Dir(conf.StrDef("www", "/www/_")))
	wapFS := http.FileServer(http.Dir(conf.StrDef("www", "/www/wap")))
	adminFS := http.FileServer(http.Dir(conf.StrDef("www", "/www/admin")))
	web.Shared.HandleNormalFunc("^.*$", func(w http.ResponseWriter, r *http.Request) {
		key := strings.SplitN(r.Host, ".", 2)[0]
		switch key {
		case "wap":
			wapFS.ServeHTTP(w, r)
		case "admin":
			adminFS.ServeHTTP(w, r)
		default:
			wwwFS.ServeHTTP(w, r)
		}
	})
	go web.HandleSignal()
	xlog.Infof("start harvester service on %v", conf.Str("/server/listen"))
	err = web.ListenAndServe(conf.Str("/server/listen"))
	if err != nil {
		panic(err)
	}
}
