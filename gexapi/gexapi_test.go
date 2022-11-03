package gexapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Centny/rediscache"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhash"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/base/sms"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/gexupgrade"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
)

const matcherConfig = `
[matcher.SPOT_YWEUSDT]
on=1
symbol=spot.YWEUSDT
base=YWE
quote=USDT
fee=0.002

[matcher.FUTURES_YWEUSDT]
on=1
symbol=futures.YWEUSDT
base=YWE
quote=USDT
fee=0.002
margin_max=0.99
margin_add=0.01
`

var ts *httptest.Server
var proxyAddr = "127.0.0.1:1105"
var ctx = context.Background()

func init() {
	func() {
		defer func() {
			recover()
		}()
		SrvAddr()
	}()
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/exservice")
	if err != nil {
		panic(err)
	}
	basedb.SYS = "exs"
	gexdb.Pool = pgx.Pool
	basedb.Pool = pgx.Pool
	gexdb.Pool().Exec(ctx, gexupgrade.DROP)
	basedb.Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "exs_"))

	//
	redisURI := "redis.loc:6379?db=1"
	rediscache.InitRedisPool(redisURI)
	gexdb.Redis = rediscache.C
	sms.Redis = rediscache.C
	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	_, err = gexdb.CheckDb(ctx)
	if err != nil {
		panic(err)
	}
	initdata()
	initMarket()
	ts = httptest.NewMuxServer()
	pgx.Client = ts.Client
	ts.Mux.HandleFunc("^/usr/mockPayTopupOrder(\\?.*)?$", MockPayTopupOrderH)
	// EnterIntentionVerifyPhoneH = NewVerifyPhone(PhoneCodeTypeVerify, "user", -1)
	Handle("", ts.Mux)
	ts.Mux.HandleNormal("^.*$", http.FileServer(http.Dir("www")))
	SrvAddr = func() string {
		return ts.URL
	}
	xhttp.EnableCookie()
	proxyServer := os.Getenv("PROXY_SERVER")
	if len(proxyServer) > 0 {
		proxyAddr = proxyServer
	}
	basedb.StoreConf(ctx, gexdb.ConfigBrokerCommRate, "0.8")
}

// func clear() {
// 	_, err := gexdb.Pool().Exec(shsupgrade.CLEAR)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func testAddUser(userRole gexdb.UserRole, account string) (user *gexdb.User) {
	phone, password, name := account+"_123", "123", account+"_name"
	image := account + "_image"
	user = &gexdb.User{
		Type:      gexdb.UserTypeNormal,
		Role:      userRole,
		Name:      &name,
		Account:   &account,
		Phone:     &phone,
		Image:     &image,
		Password:  converter.StringPtr(xhash.SHA1([]byte(password))),
		TradePass: converter.StringPtr(xhash.SHA1([]byte(password))),
		Status:    gexdb.UserStatusNormal,
	}
	err := gexdb.AddUser(ctx, user)
	if err != nil {
		panic(err)
	}
	return
}

const (
	spotBalanceBase   = "YWE"
	spotBalanceQuote  = "USDT"
	spotBalanceSymbol = "spot.YWEUSDT"
)

var spotBalanceAll = []string{spotBalanceBase, spotBalanceQuote}

var userx0, userx1 *gexdb.User
var userabc0, userabc1, userabc2, userabc3 *gexdb.User

func initdata() {
	userx0 = testAddUser(gexdb.UserRoleNormal, "x0")
	userx1 = testAddUser(gexdb.UserRoleNormal, "x1")
	userabc0 = testAddUser(gexdb.UserRoleNormal, "abc0")
	userabc1 = testAddUser(gexdb.UserRoleNormal, "abc1")

	userabc1.Status = gexdb.UserStatusLocked
	err := gexdb.UpdateUser(ctx, userabc1)
	if err != nil {
		panic(err)
	}
	userabc2 = testAddUser(gexdb.UserRoleNormal, "abc2")
	userabc3 = testAddUser(gexdb.UserRoleNormal, "abc3")
	gexdb.TouchBalance(ctx, gexdb.BalanceAreaSpot, spotBalanceAll, userx0.TID, userx1.TID, userabc0.TID, userabc1.TID, userabc2.TID, userabc3.TID)
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userabc0.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(1000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userabc0.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(1000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userabc2.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(1000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userabc2.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(1000),
		Status: gexdb.BalanceStatusNormal,
	})
}

func initMarket() {
	//bootstrap matcher
	config := xprop.NewConfig()
	config.LoadPropString(matcherConfig)
	err := matcher.Bootstrap(config)
	if err != nil {
		panic(err)
	}
	market.Bootstrap()
	//
	symbol := "spot.YWEUSDT"
	sellOpenOrder, err := matcher.ProcessLimit(ctx, userabc0.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
	if err != nil {
		panic(err)
	}
	fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)

	buyOpenOrder1, err := matcher.ProcessMarket(ctx, userabc2.TID, symbol, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
	if err != nil {
		panic(err)
	}
	fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
	buyOpenOrder2, err := matcher.ProcessLimit(ctx, userabc2.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(90))
	if err != nil {
		panic(err)
	}
	fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
	time.Sleep(300 * time.Millisecond)
}

func clearCookie() {
	xhttp.ClearCookie()
}
