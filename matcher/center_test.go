package matcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

const matcherConfig = `
[server]
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

[matcher.OFF]
on=0
`

func TestMatcherFeeCache(t *testing.T) {
	usera := testAddUser("TestMatcherFeeCache-a")
	usera.Fee = xsql.M{
		"A": 0.002,
	}
	usera.UpdateFilter(gexdb.Pool, ctx, "fee")
	userb := testAddUser("TestMatcherFeeCache-b")
	userb.Fee = xsql.M{
		"*": 0.003,
	}
	userb.UpdateFilter(gexdb.Pool, ctx, "fee")
	userc := testAddUser("TestMatcherFeeCache-c")

	cache := NewMatcherFeeCache(100)
	if fee, err := cache.LoadFee(ctx, usera.TID, "A"); err != nil || fee.InexactFloat64() != 0.002 {
		t.Errorf("%v,%v", err, fee)
		return
	}
	if fee, err := cache.LoadFee(ctx, usera.TID, "B"); err != nil || fee.InexactFloat64() != 0 {
		t.Errorf("%v,%v", err, fee)
		return
	}
	if fee, err := cache.LoadFee(ctx, userb.TID, "B"); err != nil || fee.InexactFloat64() != 0.003 {
		t.Errorf("%v,%v", err, fee)
		return
	}

	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = cache.LoadFee(ctx, userc.TID, "A")
		return
	})
}

func TestMatcherCenter(t *testing.T) {
	clear()
	config := xprop.NewConfig()
	config.LoadPropString(matcherConfig)
	center, err := BootstrapMatcherCenterByConfig(config)
	if err != nil {
		t.Error(err)
		return
	}
	center.Start()
	center.TriggerDelay = 10 * time.Millisecond
	center.eventQueue = make(chan *MatcherEvent, 1)
	eventWaiter := make(chan int, 1)
	monitor := MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
		select {
		case eventWaiter <- 1:
		default:
		}
	})
	center.AddMonitor("*", monitor)
	center.AddMonitor("spot.YWEUSDT", monitor)
	center.AddMonitor("futures.YWEUSDT", monitor)
	defer func() {
		center.RemoveMonitor("*", monitor)
	}()
	time.Sleep(time.Second)
	enabled := map[int]bool{
		0: true,
		4: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: spot buy sell cancel\n", testCount)
		//
		area := gexdb.BalanceAreaSpot
		userBase := testAddUser("TestMatcherCenter-Base")
		userQuote := testAddUser("TestMatcherCenter-Quote")
		_, err := gexdb.TouchBalance(ctx, area, spotBalanceAll, userBase.TID, userQuote.TID)
		if err != nil {
			t.Error(err)
			return
		}
		gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
			UserID: userBase.TID,
			Area:   area,
			Asset:  spotBalanceBase,
			Free:   decimal.NewFromFloat(1000),
			Status: gexdb.BalanceStatusNormal,
		})
		gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
			UserID: userQuote.TID,
			Area:   area,
			Asset:  spotBalanceQuote,
			Free:   decimal.NewFromFloat(1000),
			Status: gexdb.BalanceStatusNormal,
		})
		//
		symbol := "spot.YWEUSDT"
		sellOpenOrder, err := center.ProcessLimit(ctx, userBase.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)

		buyOpenOrder, err := center.ProcessMarket(ctx, userQuote.TID, symbol, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		cancelOrder, err := center.ProcessCancel(ctx, userBase.TID, symbol, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("cancel open order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)

		sellOpenOrderArgs := &gexdb.Order{
			UserID:       userBase.TID,
			Type:         gexdb.OrderTypeTrigger,
			OrderID:      gexdb.NewOrderID(),
			Symbol:       symbol,
			Side:         gexdb.OrderSideSell,
			Quantity:     decimal.NewFromFloat(0.5),
			Price:        decimal.NewFromFloat(100),
			TriggerType:  gexdb.OrderTriggerTypeStopProfit,
			TriggerPrice: decimal.NewFromFloat(100),
			Status:       gexdb.OrderStatusWaiting,
		}
		sellOpenOrder2, err := center.ProcessOrder(ctx, sellOpenOrderArgs) //add trigger order
		if err != nil {
			t.Error(err)
			return
		}
		_, err = center.ProcessOrder(ctx, sellOpenOrder2) //apply
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder2.OrderID)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusPending)
		cancelOrder, err = center.ProcessCancel(ctx, userBase.TID, symbol, sellOpenOrder2.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("cancel open order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusCanceled)
		<-eventWaiter

		//cache
		pgx.MockerClear()
		pgx.MockerStart()
		pgx.MockerSetCall("Pool.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = center.Preparer.PrepareSpotMatcher(ctx, center.FindMatcher(symbol).(*SpotMatcher), 100)
			return
		})
		err = center.Preparer.PrepareSpotMatcher(ctx, center.FindMatcher(symbol).(*SpotMatcher), 100)
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerSetCall("Pool.Exec", 1).Call(func(trigger int) (res xmap.M, err error) {
			err = center.Preparer.PrepareSpotMatcher(ctx, center.FindMatcher(symbol).(*SpotMatcher), 100)
			return
		})
		pgx.MockerStop()
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: futures buy sell cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		symbol := "futures.YWEUSDT"
		sellOpenOrder, err := center.ProcessLimit(ctx, env.Seller.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)

		buyOpenOrder, err := center.ProcessMarket(ctx, env.Buyer.TID, symbol, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		cancelOrder, err := center.ProcessCancel(ctx, env.Seller.TID, symbol, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("cancel open order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)
		<-eventWaiter

		//cache
		pgx.MockerClear()
		pgx.MockerStart()
		pgx.MockerSetCall("Pool.Exec", 1, "Pool.Exec", 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = center.Preparer.PrepareFuturesMatcher(ctx, center.FindMatcher(symbol).(*FuturesMatcher), 100)
			return
		})
		err = center.Preparer.PrepareFuturesMatcher(ctx, center.FindMatcher(symbol).(*FuturesMatcher), 100)
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerSetCall("Pool.Exec", 1).Call(func(trigger int) (res xmap.M, err error) {
			err = center.Preparer.PrepareFuturesMatcher(ctx, center.FindMatcher(symbol).(*FuturesMatcher), 100)
			return
		})
		pgx.MockerStop()
	}
	center.Stop()
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: trigger\n", testCount)
		//
		env := testFuturesInit(testCount)
		symbol := "futures.YWEUSDT"

		//holding
		sellOpenOrder1, err := center.ProcessLimit(ctx, env.Seller.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusPending)
		buyOpenOrder1, err := center.ProcessLimit(ctx, env.Buyer.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)

		buyCloseOrder1, err := center.ProcessOrder(ctx, &gexdb.Order{
			UserID:       env.Buyer.TID,
			Creator:      env.Buyer.TID,
			Type:         gexdb.OrderTypeTrigger,
			Symbol:       symbol,
			Side:         gexdb.OrderSideSell,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(95),
			TriggerType:  gexdb.OrderTriggerTypeStopLoss,
			TriggerPrice: decimal.NewFromFloat(95),
		})
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		center.procTriggerOrder()
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusWaiting)

		sellCloseOrder1, err := center.ProcessOrder(ctx, &gexdb.Order{
			UserID:       env.Seller.TID,
			Creator:      env.Seller.TID,
			Type:         gexdb.OrderTypeTrigger,
			Symbol:       symbol,
			Side:         gexdb.OrderSideBuy,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(95),
			TriggerType:  gexdb.OrderTriggerTypeStopProfit,
			TriggerPrice: decimal.NewFromFloat(95),
		})
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		center.procTriggerOrder()
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusWaiting)

		buyOpenOrder2, err := center.ProcessLimit(ctx, env.Buyer2.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		center.procTriggerOrder()
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)

		sellOpenOrder2, err := center.ProcessLimit(ctx, env.Seller2.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder2.OrderID)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusPending)
		center.procTriggerOrder()
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusDone)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: trigger error\n", testCount)
		//
		env := testFuturesInit(testCount)
		symbol := "futures.YWEUSDT"
		pgx.MockerStart()

		//holding
		sellOpenOrder1, err := center.ProcessLimit(ctx, env.Seller.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusPending)
		buyOpenOrder1, err := center.ProcessLimit(ctx, env.Buyer.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)

		buyCloseOrder1, err := center.ProcessOrder(ctx, &gexdb.Order{
			UserID:       env.Buyer.TID,
			Creator:      env.Buyer.TID,
			Type:         gexdb.OrderTypeTrigger,
			Symbol:       symbol,
			Side:         gexdb.OrderSideSell,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(95),
			TriggerType:  gexdb.OrderTriggerTypeStopLoss,
			TriggerPrice: decimal.NewFromFloat(95),
		})
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		center.procTriggerOrder()
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		buyOpenOrder2, err := center.ProcessLimit(ctx, env.Buyer2.TID, symbol, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		pgx.MockerClear()
		pgx.MockerSetCall("Rows.Scan", 1).Call(func(trigger int) (res xmap.M, err error) {
			err = center.procTriggerOrder()
			return
		})
		pgx.MockerPanicCall("Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = center.procTriggerOrder()
			return
		})
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		pgx.MockerClear()
		pgx.MockerSetCall("Rows.Scan", 2).Call(func(trigger int) (res xmap.M, err error) {
			err = center.procTriggerOrder()
			return
		})
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusCanceled)
		center.ProcessCancel(ctx, env.Buyer2.TID, symbol, buyOpenOrder2.OrderID)

		sellCloseOrder1, err := center.ProcessOrder(ctx, &gexdb.Order{
			UserID:       env.Seller.TID,
			Creator:      env.Seller.TID,
			Type:         gexdb.OrderTypeTrigger,
			Symbol:       symbol,
			Side:         gexdb.OrderSideBuy,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(95),
			TriggerType:  gexdb.OrderTriggerTypeStopProfit,
			TriggerPrice: decimal.NewFromFloat(95),
		})
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		center.procTriggerOrder()
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusWaiting)
		sellOpenOrder2, err := center.ProcessLimit(ctx, env.Seller2.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder2.OrderID)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusPending)
		pgx.MockerClear()
		pgx.MockerSetCall("Rows.Scan", 3).Call(func(trigger int) (res xmap.M, err error) {
			err = center.procTriggerOrder()
			return
		})
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusCanceled)
		center.ProcessCancel(ctx, env.Seller2.TID, symbol, sellOpenOrder2.OrderID)

		sellCloseOrder2, err := center.ProcessOrder(ctx, &gexdb.Order{
			UserID:       env.Seller.TID,
			Creator:      env.Seller.TID,
			Type:         gexdb.OrderTypeTrigger,
			Symbol:       symbol,
			Side:         gexdb.OrderSideBuy,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(95),
			TriggerType:  gexdb.OrderTriggerTypeStopProfit,
			TriggerPrice: decimal.NewFromFloat(95),
		})
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusWaiting)
		center.procTriggerOrder()
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusWaiting)
		sellOpenOrder3, err := center.ProcessLimit(ctx, env.Seller2.TID, symbol, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder3.OrderID)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusPending)
		pgx.MockerClear()
		pgx.MockerSetCall("Pool.Exec", 1).Call(func(trigger int) (res xmap.M, err error) {
			err = center.procTriggerOrder()
			return
		})
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusDone)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: symbol not found\n", testCount)
		//
		_, err = center.ProcessLimit(ctx, 1, "xx", gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = center.ProcessMarket(ctx, 1, "xx", gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(0.5))
		if err == nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = center.ProcessOrder(ctx, &gexdb.Order{Type: gexdb.OrderTypeTrigger, Symbol: "xx"})
		if err == nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = center.ProcessCancel(ctx, 1, "xx", "abc")
		if err == nil {
			t.Error(ErrStack(err))
			return
		}
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: error\n", testCount)
		//args eror
		_, err = center.ProcessOrder(ctx, &gexdb.Order{
			Type:   gexdb.OrderTypeBlowup,
			Symbol: "futures.YWEUSDT",
		})
		if err == nil {
			t.Error(err)
			return
		}
		_, err = center.ProcessOrder(ctx, &gexdb.Order{
			Type:   gexdb.OrderTypeTrigger,
			Symbol: "futures.YWEUSDT",
		})
		if err == nil {
			t.Error(err)
			return
		}
		_, err = center.ProcessOrder(ctx, &gexdb.Order{
			Symbol:       "futures.YWEUSDT",
			UserID:       100,
			Type:         gexdb.OrderTypeTrigger,
			Quantity:     decimal.NewFromFloat(1),
			Price:        decimal.NewFromFloat(100),
			TriggerPrice: decimal.NewFromFloat(100),
		})
		if err == nil {
			t.Error(err)
			return
		}
		center.procTriggerSybmolOrder(ctx, "xxx")
		//monitor error
		center.AddMonitor("*", MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) { panic("xxx") }))
		center.procMatcherEvent(&MatcherEvent{Symbol: "xx"})
		center.OnMatched(context.Background(), &MatcherEvent{Symbol: "xx"})
		center.OnMatched(context.Background(), &MatcherEvent{Symbol: "xx"})

		config1 := xprop.NewConfig()
		config1.LoadPropString(`
[matcher.FUTURES_YWEUSDT]
on=1
symbol=futures.YWEUSDT
		`)
		_, err = BootstrapMatcherCenterByConfig(config1)
		if err == nil {
			t.Error(err)
			return
		}

		config2 := xprop.NewConfig()
		config2.LoadPropString(`
[matcher.SPOT_YWEUSDT]
on=1
symbol=xxx.YWEUSDT
base=YWE
quote=USDT
fee=0.002
		`)
		_, err = BootstrapMatcherCenterByConfig(config2)
		if err == nil {
			t.Error(err)
			return
		}
	}
}
