package market

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
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/gexupgrade"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
	"golang.org/x/net/websocket"
)

var ctx = context.Background()

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

func init() {
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/exservice")
	if err != nil {
		panic(err)
	}
	gexdb.Pool = pgx.Pool
	basedb.SYS = "exs"
	basedb.Pool = pgx.Pool
	_, _, err = gexdb.Pool().Exec(ctx, gexupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, _, err = gexdb.Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "exs_"))
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
	err = matcher.Bootstrap(config)
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

const (
	spotBalanceBase   = "YWE"
	spotBalanceQuote  = "USDT"
	spotBalanceSymbol = "spot.YWEUSDT"
)

var spotBalanceAll = []string{spotBalanceBase, spotBalanceQuote}

func TestShared(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpot-Base")
	userQuote := testAddUser("TestSpot-Quote")
	_, err := gexdb.TouchBalance(ctx, area, append(spotBalanceAll, "NONE"), userBase.TID, userQuote.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userBase.TID,
		Area:   area,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userQuote.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	Bootstrap()
	//
	symbol := "spot.YWEUSDT"
	sellOpenOrder, err := matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)

	buyOpenOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)

	time.Sleep(300 * time.Millisecond)
	if symbols, _ := ListSymbol(); len(symbols) != 2 {
		t.Error("error")
		return
	}
	if info, line := LoadSymbol(symbol); info == nil || line == nil {
		t.Error("error")
		return
	}
	if LoadLatestPrice(symbol).Sign() <= 0 {
		t.Error("error")
		return
	}
	if len(ListLatestPrice(symbol)) < 1 {
		t.Error("error")
		return
	}
	if depth := LoadDepth(symbol, 1); depth == nil {
		t.Error("error")
		return
	}
	kline := LoadKLine(symbol, "5min")
	if kline == nil {
		t.Error("error")
		return
	}
	lines, err := ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
	if err != nil || len(lines) != 1 {
		t.Errorf("%v,%v", err, converter.JSON(lines))
		return
	}
	totalValue, areaValues, err := CalcBalanceOverview(ctx, userBase.TID)
	if err != nil || totalValue.Sign() <= 0 || len(areaValues) < 1 {
		fmt.Printf("totalValue--->%v\n", totalValue)
		fmt.Printf("totalValue--->%v\n", converter.JSON(areaValues))
		t.Error(err)
		return
	}
	pgx.MockerStart()
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, _, err = CalcBalanceOverview(ctx, userBase.TID)
		return
	})
	pgx.MockerStop()
	matcher.Shared.RemoveMonitor("*", Shared)
	Shared.Stop()

	for _, symbol := range Shared.Symbols {
		Shared.klineVal[klineKey(symbol.Symbol, "1day")] = &gexdb.KLine{
			Open:  decimal.Zero,
			Close: decimal.NewFromFloat(1),
		}
	}
	if symbols, _ := ListSymbol(); len(symbols) != 2 {
		t.Error("error")
		return
	}

	for _, symbol := range Shared.Symbols {
		Shared.klineVal[klineKey(symbol.Symbol, "1day")] = &gexdb.KLine{
			Open:  decimal.NewFromFloat(1),
			Close: decimal.NewFromFloat(1),
		}
	}
	if symbols, _ := ListSymbol(); len(symbols) != 2 {
		t.Error("error")
		return
	}
}

func TestMarketConn(t *testing.T) {
	conn := NewMarketConn(nil)
	conn.codecMarshal(xmap.M{})
	conn.codecMarshal("PING")
	conn.codecMarshal("PONG")
	conn.codecMarshal("DATA")
	conn.codecUnmarshal([]byte(``), websocket.PingFrame, nil)
	conn.codecUnmarshal([]byte(``), websocket.PongFrame, nil)
	conn.codecUnmarshal([]byte(``), websocket.TextFrame, nil)
	conn.Receive(nil)
	conn.Send(nil)
	conn.Close()
}

func TestDepthCache(t *testing.T) {
	cache := &DepthCache{}
	cache.Slice(5)
	cache.Bids = append(cache.Bids, []decimal.Decimal{})
	cache.Bids = append(cache.Bids, []decimal.Decimal{})
	cache.Asks = append(cache.Asks, []decimal.Decimal{})
	cache.Asks = append(cache.Asks, []decimal.Decimal{})
	cache.Slice(1)
}

func TestMarket(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpot-Base")
	userQuote := testAddUser("TestSpot-Quote")
	_, err := gexdb.TouchBalance(ctx, area, append(spotBalanceAll, "NONE"), userBase.TID, userQuote.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userBase.TID,
		Area:   area,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userQuote.TID,
		Area:   area,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})

	enabled := map[int]bool{
		0: true,
		1: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: market notify\n", testCount)

		market := NewMarket(matcher.Shared.Symbols)
		market.Symbols = matcher.Shared.Symbols
		market.OnConnect = func(conn *websocket.Conn) {}
		market.OnDisconnect = func(conn *websocket.Conn) {}
		market.KLineNotifyDelay = 10 * time.Millisecond
		market.KLineGenDelay = 10 * time.Millisecond
		market.WaitTimeout = 300 * time.Millisecond
		matcher.Shared.AddMonitor("*", market)
		ts := httptest.NewMuxServer()
		ts.Mux.Handle("/ws", market)
		market.Start()
		symbol := "spot.YWEUSDT"
		waiter := sync.WaitGroup{}

		waiter.Add(1)
		go func() { //connect error
			defer waiter.Done()
			conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/ws", "", ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			conn.Write([]byte(converter.JSON(xmap.M{}))) //ping
			conn.Write([]byte(converter.JSON(xmap.M{     //sub fail
				"action": "sub.kline",
			})))
			conn.Write([]byte(converter.JSON(xmap.M{ //err command
				"action": "xxx",
			})))
			conn.Close()
		}()
		waiter.Add(1)
		go func() { //wait timeout
			defer waiter.Done()
			conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/ws", "", ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			buff := make([]byte, 4096)
			conn.Read(buff)
			conn.Close()
		}()
		waiter.Add(1)
		go func() { //kline
			defer waiter.Done()
			conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/ws", "", ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			conn.Write([]byte(converter.JSON(xmap.M{
				"action": "sub.kline",
				"symbols": []xmap.M{
					{
						"symbol":   symbol,
						"interval": "5min",
					},
				},
			})))
			buff := make([]byte, 4096)
			received := 0
			for {
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Printf("receive:%v\n", string(buff[0:n]))
				received++
				if received > 5 {
					break
				}
			}
			conn.Close()
		}()
		waiter.Add(1)
		go func() { //depth
			defer waiter.Done()
			conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/ws", "", ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			conn.Write([]byte(converter.JSON(xmap.M{
				"action": "sub.depth",
				"symbols": []xmap.M{
					{
						"symbol": symbol,
						"max":    1,
					},
				},
			})))
			buff := make([]byte, 4096)
			received := 0
			for {
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Printf("receive:%v\n", string(buff[0:n]))
				received++
				if received > 5 {
					break
				}
			}
			conn.Close()
		}()
		waiter.Add(1)
		go func() { //depth
			defer waiter.Done()
			conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/ws", "", ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			conn.Write([]byte(converter.JSON(xmap.M{
				"action":  "sub.ticker",
				"symbols": []string{symbol},
			})))
			buff := make([]byte, 4096)
			received := 0
			for {
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Printf("receive:%v\n", string(buff[0:n]))
				received++
				if received > 5 {
					break
				}
			}
			conn.Close()
		}()
		time.Sleep(500 * time.Millisecond)

		_, err = matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(err)
			return
		}
		waiter.Wait()
		market.Stop()
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: proc kline\n", testCount)
		//
		market := NewMarket(matcher.Shared.Symbols)
		market.Symbols = matcher.Shared.Symbols
		market.OnConnect = func(conn *websocket.Conn) {}
		market.OnDisconnect = func(conn *websocket.Conn) {}
		market.KLineNotifyDelay = 10 * time.Millisecond
		market.KLineGenDelay = 10 * time.Millisecond
		market.WaitTimeout = 300 * time.Millisecond
		symbol := "spot.YWEUSDT"

		//first line
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
			Orders: []*gexdb.Order{
				{
					Filled:     decimal.NewFromFloat(1),
					AvgPrice:   decimal.NewFromFloat(100),
					TotalPrice: decimal.NewFromFloat(100),
				},
			},
		})
		kline := market.LoadKLine(symbol, "5min")
		if kline == nil {
			t.Error("error")
			return
		}
		lines, err := market.ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
		if err != nil || len(lines) != 1 {
			t.Errorf("%v,%v", err, converter.JSON(lines))
			return
		}

		//save kline
		for _, line := range market.klineVal {
			interval, _ := gexdb.StringInterv(line.Interv)
			line.StartTime = xsql.Time(line.StartTime.AsTime().Add(-interval))
		}
		market.procGenKLine(nil)
		lines, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
		if err != nil || len(lines) != 2 {
			t.Errorf("%v,%v", err, converter.JSON(lines))
			return
		}

		//new line
		for _, line := range market.klineVal {
			interval, _ := gexdb.StringInterv(line.Interv)
			line.StartTime = xsql.Time(line.StartTime.AsTime().Add(-interval))
		}
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
			Orders: []*gexdb.Order{
				{
					Filled:     decimal.NewFromFloat(1),
					AvgPrice:   decimal.NewFromFloat(100),
					TotalPrice: decimal.NewFromFloat(100),
				},
			},
		})
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
			Orders: []*gexdb.Order{
				{
					Filled:     decimal.NewFromFloat(1),
					AvgPrice:   decimal.NewFromFloat(110),
					TotalPrice: decimal.NewFromFloat(110),
				},
			},
		})
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
			Orders: []*gexdb.Order{
				{
					Filled:     decimal.NewFromFloat(1),
					AvgPrice:   decimal.NewFromFloat(90),
					TotalPrice: decimal.NewFromFloat(90),
				},
			},
		})
		lines, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
		if err != nil || len(lines) != 3 {
			t.Errorf("%v,%v", err, converter.JSON(lines))
			return
		}
		lines, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
		if err != nil || len(lines) != 5 { //load more from db
			t.Errorf("%v,%v", err, converter.JSON(lines))
			return
		}
		pgx.MockerSetCall("Pool.Query", 1).Call(func(trigger int) (res xmap.M, err error) { //should not error, using cache
			_, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-1000*time.Minute), time.Now())
			return
		})

		//error
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
		})
		pgx.MockerClear()
		for _, line := range market.klineVal {
			interval, _ := gexdb.StringInterv(line.Interv)
			line.StartTime = xsql.Time(line.StartTime.AsTime().Add(-interval))
		}
		pgx.MockerSetCall("Pool.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = market.procGenKLine(nil)
			return
		})
		for _, line := range market.klineVal {
			interval, _ := gexdb.StringInterv(line.Interv)
			line.StartTime = xsql.Time(line.StartTime.AsTime().Add(-interval))
		}
		pgx.MockerPanicCall("Pool.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = market.procGenKLine(nil)
			return
		})

		//list error
		_, err = market.ListKLine(ctx, symbol, "xxxx", time.Now().Add(-5*time.Minute), time.Now())
		if err == nil {
			t.Error(err)
			return
		}
		delete(market.klineVal, klineKey(symbol, "5min"))
		pgx.MockerClear()
		pgx.MockerSetCall("Pool.Query", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-5*time.Minute), time.Now())
			return
		})
		_, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-5*time.Minute), time.Now())
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()
		pgx.MockerSetCall("Pool.Query", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) { //more data
			_, err = market.ListKLine(ctx, symbol, "5min", time.Now().Add(-15*time.Minute), time.Now())
			return
		})
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: proc error\n", testCount)
		//
		//
		market := NewMarket(matcher.Shared.Symbols)
		market.Symbols = matcher.Shared.Symbols
		market.OnConnect = func(conn *websocket.Conn) {}
		market.OnDisconnect = func(conn *websocket.Conn) {}
		market.KLineNotifyDelay = 10 * time.Millisecond
		market.KLineGenDelay = 10 * time.Millisecond
		market.WaitTimeout = 300 * time.Millisecond
		symbol := "spot.YWEUSDT"

		//notfy kline
		market.procGenKLine(&matcher.MatcherEvent{
			Symbol: symbol,
			Orders: []*gexdb.Order{
				{
					Filled:     decimal.NewFromFloat(1),
					AvgPrice:   decimal.NewFromFloat(100),
					TotalPrice: decimal.NewFromFloat(100),
				},
			},
		})
		market.procNotifyKLine(&MarketConn{
			KLines: map[string]int{
				klineKey(symbol, "5min"): 1,
			},
		})
		market.procNotifyKLine(&MarketConn{
			Conn: &websocket.Conn{},
			KLines: map[string]int{
				klineKey(symbol, "5min"): 1,
			},
		})

		//notify depth
		depth := &DepthCache{
			Symbol: symbol,
			Time:   xsql.TimeNow(),
		}
		market.depthVal[depth.Symbol] = depth
		market.procNotifyDepth(&MarketConn{
			Depths: map[string]int{symbol: 1},
		}, "xxx")
		market.procNotifyDepth(&MarketConn{
			Depths: map[string]int{"xxx": 1},
		}, "xxx")
		market.procNotifyDepth(&MarketConn{
			Depths: map[string]int{symbol: 1},
		}, depth.Symbol)
		market.procNotifyDepth(&MarketConn{
			Conn:   &websocket.Conn{},
			Depths: map[string]int{symbol: 1},
		}, depth.Symbol)

		//notify ticker
		depth = &DepthCache{
			Symbol: symbol,
			Time:   xsql.TimeNow(),
			Asks:   [][]decimal.Decimal{{decimal.NewFromFloat(2), decimal.NewFromFloat(1)}},
			Bids:   [][]decimal.Decimal{{decimal.NewFromFloat(1), decimal.NewFromFloat(1)}},
		}
		market.depthVal[depth.Symbol] = depth
		market.depthVal[depth.Symbol] = depth
		market.procNotifyTicker(&MarketConn{
			Tickers: map[string]int{symbol: 1},
		}, depth.Symbol)
		market.procNotifyTicker(&MarketConn{
			Conn:    &websocket.Conn{},
			Tickers: map[string]int{symbol: 1},
		}, depth.Symbol)

		//queue full
		market.eventQueue = make(chan *matcher.MatcherEvent, 1)
		market.eventQueue <- nil
		market.OnMatched(ctx, &matcher.MatcherEvent{})
		market.depthQueue = make(chan *depthQueueItem, 1)
		market.depthQueue <- nil
		market.depthVal[depth.Symbol] = depth
		market.wsconn["xxx"] = &MarketConn{
			Depths: map[string]int{symbol: 1},
		}
		market.procTriggerDepth(&matcher.MatcherEvent{Symbol: symbol, Depth: &orderbook.Depth{}})

		//trigger kline
		var nilMarket *Market
		nilMarket.procTriggerKLine()
		nilMarket.procTriggerDepth(nil)
	}

	// time.Sleep(100 * time.Millisecond)
	// market.Stop()

	// //test error
	// pgx.MockerStart()
	// defer pgx.MockerStop()
	// pgx.MockerClear()

	// pgx.MockerSet("Pool.Exec", 1)
	// for _, line := range market.klineVal {
	// 	interval, _ := StringInterv(line.Interv)
	// 	line.StartTime = xsql.Time(line.StartTime.AsTime().Add(-interval))
	// }
	// err = market.procKLine(&MatcherEvent{
	// 	Order: &Order{
	// 		Filled:     decimal.NewFromFloat(1),
	// 		AvgPrice:   decimal.NewFromFloat(100),
	// 		TotalPrice: decimal.NewFromFloat(100),
	// 	},
	// })
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// market.wsconn["x"] = &MarketConn{}
	// market.wsconn["y"] = &MarketConn{Latest: time.Now(), KLineInterv: []string{"5min"}}
	// market.procNotify()

	// market.eventQueue = make(chan *MatcherEvent)
	// market.OnMatched(&MatcherEvent{})

	// var nilm *Market
	// nilm.procKLine(nil)
	// nilm.procNotify()

	// //
	// _, err = market.ListKLine(ctx, "xxx", time.Now(), time.Now())
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// pgx.MockerSet("Pool.Query", 1)
	// _, err = market.ListKLine(ctx, "30min", time.Now().Add(-30*time.Minute), time.Now())
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// _, err = market.ListKLine(ctx, "30min", time.Now().Add(-30*time.Minute), time.Now())
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// pgx.MockerSet("Pool.Query", 1)
	// _, err = market.ListKLine(ctx, "30min", time.Now().Add(-60*time.Minute), time.Now())
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// market.LatestPrice()
}
