package maker

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/gexupgrade"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
)

var ctx = context.Background()

const matcherConfig = `
[matcher.SPOT_YWEUSDT]
on=1
type=spot
symbol=spot.YWEUSDT
base=YWE
quote=USDT
precision_quantity=2
precision_price=2
fee=0.002

[matcher.FUTURES_YWEUSDT]
on=1
type=futures
symbol=futures.YWEUSDT
base=YWE
quote=USDT
precision_quantity=2
precision_price=2
fee=0.002
margin_max=0.99
margin_add=0.01
`

func init() {
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/gexservice")
	if err != nil {
		panic(err)
	}
	gexdb.Pool = pgx.Pool
	basedb.SYS = "gex"
	basedb.Pool = pgx.Pool
	_, _, err = gexdb.Pool().Exec(ctx, gexupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, _, err = gexdb.Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "gex_"))
	if err != nil {
		panic(err)
	}
	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	_, err = gexdb.CheckDb(ctx)
	if err != nil {
		panic(err)
	}
	config := xprop.NewConfig()
	config.LoadPropString(matcherConfig)
	err = matcher.Bootstrap(ctx, config)
	if err != nil {
		panic(err)
	}
	go http.ListenAndServe(":6062", nil)
}

func clear() {
	_, _, err := gexdb.Pool().Exec(ctx, gexupgrade.CLEAR)
	if err != nil {
		panic(err)
	}
}

func testAddUser(prefix string) (user *gexdb.User) {
	account, phone, password := prefix+"_acc", prefix+"_123", "123"
	image := prefix + "_image"
	user = &gexdb.User{
		Type:      gexdb.UserTypeNormal,
		Role:      gexdb.UserRoleNormal,
		Name:      &prefix,
		Account:   &account,
		Phone:     &phone,
		Image:     &image,
		Password:  &password,
		TradePass: &password,
		External:  xsql.M{"abc": 1},
		Status:    gexdb.UserStatusNormal,
	}
	err := gexdb.AddUser(ctx, user)
	if err != nil {
		panic(err)
	}
	return
}

func TestControl(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	area := gexdb.BalanceAreaSpot
	userMaker := testAddUser("TestMakerSpot-Maker")
	userTaker := testAddUser("TestMakerSpot-Taker")
	_, err := gexdb.TouchBalance(ctx, area, append(spotBalanceAll, "NONE"), userMaker.TID, userTaker.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userMaker.TID,
		Area:   area,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userMaker.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userTaker.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	config := Config{}
	err = config.Valid()
	if err == nil {
		t.Error(err)
		return
	}
	config.ON = 1
	config.Symbol = spotSymbol
	config.Delay = 500
	config.UserID = userMaker.TID
	config.Open = decimal.NewFromFloat(1000)
	config.Close.Min = decimal.NewFromFloat(-0.01)
	config.Close.Max = decimal.NewFromFloat(0.01)
	config.Vib.Min = decimal.NewFromFloat(-0.03)
	config.Vib.Max = decimal.NewFromFloat(0.03)
	config.Vib.Count = 5
	config.Ticker = decimal.NewFromFloat(0.0001)
	config.Interval = time.Hour.Milliseconds()
	config.Depth.QtyMax = decimal.NewFromFloat(3)
	config.Depth.StepMax = 5
	config.Depth.DiffMax = decimal.NewFromFloat(2)
	config.Depth.DiffMin = decimal.NewFromFloat(0.02)
	config.Depth.Max = 15
	err = config.Valid()
	if err != nil {
		t.Error(err)
		return
	}
	err = basedb.StoreConf(ctx, fmt.Sprintf("maker-%v", config.Symbol), converter.JSON(config))
	if err != nil {
		t.Error(err)
		return
	}
	err = Bootstrap(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	err = Bootstrap(ctx)
	if err == nil {
		t.Error(err)
		return
	}
	err = Start(ctx, config.Symbol)
	if err == nil {
		t.Error(err)
		return
	}
	maker := Find(ctx, config.Symbol)
	if maker == nil {
		t.Error(err)
		return
	}
	makers, _ := List(ctx)
	if len(makers) < 1 {
		t.Error(err)
		return
	}
	err = UpdateConfig(ctx, &config)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = LoadConfig(ctx, config.Symbol)
	if err != nil {
		t.Error(err)
		return
	}
	err = Stop(ctx, config.Symbol)
	if err != nil {
		t.Error(err)
		return
	}
	err = Stop(ctx, config.Symbol)
	if err == nil {
		t.Error(err)
		return
	}
	maker.Config.Symbol = "xx"
	err = maker.Start(ctx)
	if err == nil {
		t.Error(err)
		return
	}

	pgx.MockerClear()

	pgx.MockerSetCall("Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = Bootstrap(ctx)
		return
	})

	pgx.MockerSetCall("Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = Start(ctx, config.Symbol)
		return
	})

	pgx.MockerSetCall("Pool.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = UpdateConfig(ctx, &config)
		return
	})
}

func TestConfig(t *testing.T) {
	config := &Config{}
	config.Open = decimal.NewFromFloat(1000)
	config.Close.Min = decimal.NewFromFloat(-0.01)
	config.Close.Max = decimal.NewFromFloat(0.01)
	config.Vib.Min = decimal.NewFromFloat(-0.03)
	config.Vib.Max = decimal.NewFromFloat(0.03)
	config.Vib.Count = 5
	config.Ticker = decimal.NewFromFloat(0.0001)
	config.Interval = time.Hour.Milliseconds()
	past := time.Duration(0)
	close := config.Open
	for {
		next := config.Random(past, close)
		// fmt.Printf("%v\n", next.InexactFloat64())
		close = next
		past += time.Second
		if past.Milliseconds() > config.Interval {
			break
		}
	}
	config.Value()
	config.Scan(t)
	config = nil
	config.Value()
}

const (
	spotBalanceBase     = "YWE"
	spotBalanceQuote    = "USDT"
	spotSymbol          = "spot.YWEUSDT"
	futuresBalanceQuote = "USDT"
	futuresSymbol       = "futures.YWEUSDT"
)

var spotBalanceAll = []string{spotBalanceBase, spotBalanceQuote}
var futuresBalanceAll = []string{futuresBalanceQuote}

func TestMakerSpot(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	area := gexdb.BalanceAreaSpot
	userMaker := testAddUser("TestMakerSpot-Maker")
	userTaker := testAddUser("TestMakerSpot-Taker")
	_, err := gexdb.TouchBalance(ctx, area, append(spotBalanceAll, "NONE"), userMaker.TID, userTaker.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userMaker.TID,
		Area:   area,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userMaker.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userTaker.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	config := Config{}
	config.Symbol = spotSymbol
	config.Delay = 20
	config.UserID = userMaker.TID
	config.Open = decimal.NewFromFloat(1000)
	config.Close.Min = decimal.NewFromFloat(-0.01)
	config.Close.Max = decimal.NewFromFloat(0.01)
	config.Vib.Min = decimal.NewFromFloat(-0.03)
	config.Vib.Max = decimal.NewFromFloat(0.03)
	config.Vib.Count = 5
	config.Ticker = decimal.NewFromFloat(0.0001)
	config.Interval = time.Hour.Milliseconds()
	config.Depth.QtyMax = decimal.NewFromFloat(3)
	config.Depth.StepMax = 5
	config.Depth.DiffMax = decimal.NewFromFloat(2)
	config.Depth.DiffMin = decimal.NewFromFloat(0.02)
	config.Depth.Max = 15
	maker := NewMaker(&config)
	{ //test willPlace
		maker.symbol = matcher.Shared.Symbols[config.Symbol]
		maker.balances = map[string]*gexdb.Balance{}
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{}
		maker.makingAll = map[string]decimal.Decimal{}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{Asset: maker.symbol.Quote, Free: decimal.NewFromFloat(101)}
		if !maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.balances[maker.symbol.Base] = &gexdb.Balance{Asset: maker.symbol.Quote, Free: decimal.NewFromFloat(1)}
		if !maker.willPlace(gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.makingAll["100"] = decimal.NewFromFloat(100000)
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{Asset: maker.symbol.Quote, Free: decimal.NewFromFloat(1000)}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
	}
	maker.Verbose = true
	matcher.Shared.AddMonitor(config.Symbol, maker)
	err = maker.Start(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	maker.clearShow = time.Time{}
	maker.nextShow = time.Time{}
	time.Sleep(1000 * time.Millisecond)
	order, err := matcher.ProcessMarket(ctx, userTaker.TID, config.Symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(10), decimal.Zero)
	if err != nil || order.Status == gexdb.OrderStatusCanceled {
		t.Error(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	maker.Update(&config)
	maker.Stop()
	fmt.Printf("--->%v\n", maker.next)
	fmt.Printf("--->%v\n", converter.JSON(maker.depth))

	//error
	maker.balances[spotBalanceQuote].Free = decimal.NewFromFloat(100000000000000)
	maker.procPlace(ctx, gexdb.OrderSideBuy, decimal.NewFromFloat(10000000))

	maker.makingOrder["xxx"] = &gexdb.Order{OrderID: "xxx", Side: gexdb.OrderSideBuy, Price: decimal.NewFromFloat(10)}
	maker.procCancle(ctx, decimal.NewFromFloat(9), decimal.NewFromFloat(9))

	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Exec", 1, "Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = maker.Start(ctx)
		return
	})

	//
	maker.makerQueue = make(chan int, 1)
	maker.depth = &orderbook.Depth{Asks: [][]decimal.Decimal{{decimal.NewFromFloat(3), decimal.NewFromFloat(1)}}, Bids: [][]decimal.Decimal{{decimal.NewFromFloat(2), decimal.NewFromFloat(1)}}}
	maker.OnMatched(ctx, &matcher.MatcherEvent{
		Depth: &orderbook.Depth{Asks: [][]decimal.Decimal{{decimal.NewFromFloat(2), decimal.NewFromFloat(1)}}, Bids: [][]decimal.Decimal{{decimal.NewFromFloat(1), decimal.NewFromFloat(1)}}},
	})
	//
	maker.makingOrder["xx"] = &gexdb.Order{Status: gexdb.OrderStatusCanceled}
	maker.checkOrder(decimal.NewFromFloat(1), decimal.NewFromFloat(1))
	//
	maker.clearLast = time.Time{}
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Exec", 1).Call(func(trigger int) (res xmap.M, err error) {
		maker.procClear(ctx)
		return
	})
}

func TestMakerFutures(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	area := gexdb.BalanceAreaFutures
	userMaker := testAddUser("TestMakerFutures-Maker")
	userTaker := testAddUser("TestMakerFutures-Taker")
	_, err := gexdb.TouchBalance(ctx, area, futuresBalanceAll, userMaker.TID, userTaker.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userMaker.TID,
		Area:   area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(100000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userTaker.TID,
		Area:   area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(100000),
		Status: gexdb.BalanceStatusNormal,
	})
	config := Config{}
	config.Symbol = futuresSymbol
	config.Delay = 20
	config.UserID = userMaker.TID
	config.Open = decimal.NewFromFloat(1000)
	config.Close.Min = decimal.NewFromFloat(-0.01)
	config.Close.Max = decimal.NewFromFloat(0.01)
	config.Vib.Min = decimal.NewFromFloat(-0.03)
	config.Vib.Max = decimal.NewFromFloat(0.03)
	config.Vib.Count = 5
	config.Ticker = decimal.NewFromFloat(0.0001)
	config.Interval = time.Hour.Milliseconds()
	config.Depth.QtyMax = decimal.NewFromFloat(1)
	config.Depth.StepMax = 5
	config.Depth.DiffMax = decimal.NewFromFloat(2)
	config.Depth.DiffMin = decimal.NewFromFloat(0.02)
	config.Depth.Max = 15
	maker := NewMaker(&config)
	{ //test willPlace
		maker.holding = nil
		maker.symbol = matcher.Shared.Symbols[config.Symbol]
		maker.balances = map[string]*gexdb.Balance{}
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{}
		maker.makingAll = map[string]decimal.Decimal{}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{Asset: maker.symbol.Quote, Free: decimal.NewFromFloat(101)}
		if !maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.holding = &gexdb.Holding{Amount: decimal.NewFromFloat(1)}
		if !maker.willPlace(gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.holding = &gexdb.Holding{Amount: decimal.NewFromFloat(-1)}
		if !maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
		maker.makingAll["100"] = decimal.NewFromFloat(100000)
		maker.balances[maker.symbol.Quote] = &gexdb.Balance{Asset: maker.symbol.Quote, Free: decimal.NewFromFloat(1000)}
		if maker.willPlace(gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100)) {
			t.Error("error")
			return
		}
	}
	maker.Verbose = true
	matcher.Shared.AddMonitor(config.Symbol, maker)
	err = maker.Start(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	maker.clearShow = time.Time{}
	maker.nextShow = time.Time{}
	time.Sleep(500 * time.Millisecond)
	waiter := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			order, err := matcher.ProcessMarket(ctx, userTaker.TID, config.Symbol, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(100))
			if err != nil || order.TID < 1 {
				t.Error(err)
				return
			}
		}()
	}
	waiter.Wait()
	time.Sleep(1000 * time.Millisecond)
	maker.Stop()
	fmt.Printf("--->%v\n", maker.next)
	fmt.Printf("--->%v\n", converter.JSON(maker.depth))

	//error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Exec", 1, 2, "Rows.Scan", 1, 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = maker.Start(ctx)
		return
	})
}
