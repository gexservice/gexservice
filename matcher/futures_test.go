package matcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

const (
	futuresBalanceQuote  = "USDT"
	futuresHoldingSymbol = "YWEUSDT"
)

var futuresBalanceAll = []string{futuresBalanceQuote}
var futuresHoldingAll = []string{futuresHoldingSymbol}

type FuturesTestEnv struct {
	Area           gexdb.BalanceArea
	Buyer          *gexdb.User
	BuyerHolding   *gexdb.Holding
	Buyer2         *gexdb.User
	BuyerHolding2  *gexdb.Holding
	Seller         *gexdb.User
	SellerHolding  *gexdb.Holding
	Seller2        *gexdb.User
	SellerHolding2 *gexdb.Holding
	Adder          *gexdb.User
	AdderHolding   *gexdb.Holding
	Small          *gexdb.User
	SmallHolding   *gexdb.Holding
	None           *gexdb.User
	NoneHolding    *gexdb.Holding
	Monitor        MatcherMonitor
}

func testFuturesInit(i int) (env *FuturesTestEnv) {
	env = &FuturesTestEnv{}
	env.Area = gexdb.BalanceAreaFutures
	env.Buyer = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Buy-%v", i))
	env.Seller = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Sell-%v", i))
	env.Buyer2 = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Buy2-%v", i))
	env.Seller2 = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Sell2-%v", i))
	env.Adder = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Add-%v", i))
	env.Small = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-Small-%v", i))
	env.None = testAddUser(fmt.Sprintf("TestFuturesMatcherMarket-NONE-%v", i))
	_, err := gexdb.TouchBalance(ctx, env.Area, futuresBalanceAll, env.Buyer.TID, env.Seller.TID, env.Buyer2.TID, env.Seller2.TID, env.Adder.TID, env.Small.TID, env.None.TID)
	if err != nil {
		panic(err)
	}
	_, err = gexdb.TouchHolding(ctx, futuresHoldingAll, env.Buyer.TID, env.Seller.TID, env.Buyer2.TID, env.Seller2.TID, env.Adder.TID, env.Small.TID, env.None.TID)
	if err != nil {
		panic(err)
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Buyer.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Seller.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Buyer2.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Seller2.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Adder.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: env.Small.TID,
		Area:   env.Area,
		Asset:  futuresBalanceQuote,
		Free:   decimal.NewFromFloat(21),
		Status: gexdb.BalanceStatusNormal,
	})

	env.BuyerHolding, err = gexdb.FindHoldlingBySymbol(ctx, env.Buyer.TID, futuresHoldingSymbol)
	if err != nil {
		panic(err)
	}
	env.BuyerHolding.Lever = 10
	err = gexdb.UpdateHoldingFilter(ctx, env.BuyerHolding, "lever")
	if err != nil {
		panic(err)
	}

	env.SellerHolding, err = gexdb.FindHoldlingBySymbol(ctx, env.Seller.TID, futuresHoldingSymbol)
	if err != nil {
		panic(err)
	}
	env.SellerHolding.Lever = 5
	err = gexdb.UpdateHoldingFilter(ctx, env.SellerHolding, "lever")
	if err != nil {
		panic(err)
	}

	env.AdderHolding, err = gexdb.FindHoldlingBySymbol(ctx, env.Adder.TID, futuresHoldingSymbol)
	if err != nil {
		panic(err)
	}
	env.AdderHolding.Lever = 1
	err = gexdb.UpdateHoldingFilter(ctx, env.AdderHolding, "lever")
	if err != nil {
		panic(err)
	}

	env.SmallHolding, err = gexdb.FindHoldlingBySymbol(ctx, env.Small.TID, futuresHoldingSymbol)
	if err != nil {
		panic(err)
	}
	env.SmallHolding.Lever = 10
	err = gexdb.UpdateHoldingFilter(ctx, env.SmallHolding, "lever")
	if err != nil {
		panic(err)
	}

	env.NoneHolding, err = gexdb.FindHoldlingBySymbol(ctx, env.None.TID, futuresHoldingSymbol)
	if err != nil {
		panic(err)
	}
	env.Monitor = MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {})
	return
}

func TestFuturesMatcherBootstrap(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		3: true,
	}
	pgx.MockerStart()
	defer pgx.MockerStop()
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		changed, err := matcher.Bootstrap(ctx)
		if err != nil || len(changed.Orders) > 0 {
			t.Error(ErrStack(err))
			return
		}
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		pgx.MockerClear()
		matcher = NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		pgx.MockerPanicCall("Pool.Begin", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.Bootstrap(ctx)
			return
		})
		pgx.MockerSetCall("Pool.Begin", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.Bootstrap(ctx)
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.Bootstrap(ctx)
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.Bootstrap(ctx)
			return
		})
		changed, err = matcher.Bootstrap(ctx)
		if err != nil || len(changed.Orders) < 1 {
			t.Error(err)
			return
		}
	}
}

func TestFuturesMatcherMarket(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		2: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all, invest\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(110), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all, invest, prepare\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder := &gexdb.Order{
			Type:       gexdb.OrderTypeTrigger,
			OrderID:    matcher.NewOrderID(),
			UserID:     env.Seller.TID,
			Symbol:     matcher.Symbol,
			Side:       gexdb.OrderSideBuy,
			TotalPrice: decimal.NewFromFloat(110),
			Status:     gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all, quantity\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy partial, invest\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder1, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(50), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		buyOpenOrder2, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(50), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder1.OrderID)
		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder2.OrderID)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder1, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder1.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusDone)

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder2, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(10), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder2.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy partial, quantity\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder1, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		buyOpenOrder2, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.5))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder1.OrderID)
		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder2.OrderID)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder1, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.8))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder1.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusDone)

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder2, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.2))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder2.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy part cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(200), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(110), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy close and open, invest\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		fmt.Printf("%v start sell open order\n", env.Adder.TID)
		addOpenOrder, err := matcher.ProcessLimit(ctx, env.Adder.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell open order %v\n", env.Adder.TID, addOpenOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(addOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close=>open order\n", env.Seller.TID)
		sellCloseOpenOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(220), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close=>open order %v\n", env.Seller.TID, sellCloseOpenOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(addOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Adder.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(110))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(22))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Adder.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(110))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(22))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Adder.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9889.78))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9967.36))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Adder.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start buy close order\n", env.Adder.TID)
		addCloseOrder, err := matcher.ProcessLimit(ctx, env.Adder.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Adder.TID, addCloseOrder.OrderID)
		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.14))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell multi buy all, quantity\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder1, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		sellOpenOrder2, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder2.OrderID)
		sellOpenOrder3, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(200))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder3.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(60.6))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(2))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(30))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(60))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(30))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(60))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9969.4))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9939.4))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(2))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-2))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(150))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(30))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(30.6))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(2))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(2))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9998.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9998.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: buy cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)

		buyOpenOrder1, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(200), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)

		buyOpenOrder2, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusCanceled)

		sellOpenOrder1, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	//test error
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: buy not enought, invest\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		_, err = matcher.ProcessMarket(ctx, env.None.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if !IsErrBalanceNotEnought(err) {
			t.Error(err)
			return
		}
		fmt.Printf("process err is \n%v\n", ErrStack(err))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: arg err\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		//
		_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, 0, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, env.None.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, env.None.TID, gexdb.OrderSideSell, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, env.None.TID, "0", decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}

		//
		doneOrdeer := &gexdb.Order{
			Type:       gexdb.OrderTypeTrigger,
			OrderID:    matcher.NewOrderID(),
			UserID:     env.Seller.TID,
			Symbol:     matcher.Symbol,
			Side:       gexdb.OrderSideBuy,
			TotalPrice: decimal.NewFromFloat(110),
			Status:     gexdb.OrderStatusDone,
		}
		err = gexdb.AddOrder(ctx, doneOrdeer)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = matcher.ProcessOrder(ctx, doneOrdeer)
		if err == nil {
			t.Error(ErrStack(err))
			return
		}

		//prepare error
		matcher.PrepareProcess = func(ctx context.Context, matcher *FuturesMatcher, userID int64) error {
			return fmt.Errorf("error")
		}
		_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
	}
}

func TestFuturesMatcherLimit(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		6: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy all, preorder\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder := &gexdb.Order{
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			UserID:   env.Buyer.TID,
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideSell,
			Quantity: decimal.NewFromFloat(1),
			Price:    decimal.NewFromFloat(110),
			Status:   gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, buyCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = matcher.ProcessOrder(ctx, buyCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder := &gexdb.Order{
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			UserID:   env.Seller.TID,
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideBuy,
			Quantity: decimal.NewFromFloat(1),
			Price:    decimal.NewFromFloat(110),
			Status:   gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: buy sell, buy sell partial\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.1))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-0.5))

		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Buyer.TID)
		buyCloseOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Buyer.TID, buyCloseOrder1.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.11))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusPartialled)

		fmt.Printf("%v start buy close order\n", env.Buyer.TID)
		buyCloseOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Buyer.TID, buyCloseOrder2.OrderID)
		assetOrderStatus(buyCloseOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: holding close open\n", testCount)
		//
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder1, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start sell close order\n", env.Small.TID)
		sellCloseOrder1, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Small.TID, sellCloseOrder1.OrderID)
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.1))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start sell close order\n", env.Small.TID)
		sellCloseOrder2, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Small.TID, sellCloseOrder2.OrderID)
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.4))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start sell close order\n", env.Small.TID)
		sellCloseOrder3, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.04), decimal.NewFromFloat(100))
		if !IsErrBalanceNotEnought(err) || sellCloseOrder3.Status != gexdb.OrderStatusCanceled {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Small.TID, sellCloseOrder3.OrderID)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.4))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy open order\n", env.Buyer.TID)
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Buyer.TID, buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder2.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.4))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.4))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
	}
	//test error
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: not enought\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		_, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1000000), decimal.NewFromFloat(100))
		if !IsErrBalanceNotEnought(err) {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(100000000), decimal.NewFromFloat(100))
		if !IsErrBalanceNotEnought(err) {
			t.Error(err)
			return
		}
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: arg err\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)

		//arg error
		_, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, "0", decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}

		//status error
		doneOrder := &gexdb.Order{
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			UserID:   env.Buyer.TID,
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideSell,
			Quantity: decimal.NewFromFloat(1),
			Price:    decimal.NewFromFloat(110),
			Status:   gexdb.OrderStatusDone,
		}
		err = gexdb.AddOrder(ctx, doneOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		_, err = matcher.ProcessOrder(ctx, doneOrder)
		if err == nil {
			t.Error(ErrStack(err))
			return
		}

		//prepare error
		matcher.PrepareProcess = func(ctx context.Context, matcher *FuturesMatcher, userID int64) error {
			return fmt.Errorf("error")
		}
		_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
	}
}

func TestFuturesMatcherCancel(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		6: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: buy cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)

		cancelOrder2, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
		if !IsErrNotCancelable(err) {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder2.OrderID)
		assetOrderStatus(cancelOrder2.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy sell part cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(40.4))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(40.2))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9959.6))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: sell buy buy part cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.6))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: close cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyCloseOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: cancel not access\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Seller.TID, buyOpenOrder.OrderID)
		if err != define.ErrNotAccess {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: arg err\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		_, err := matcher.ProcessCancel(ctx, env.Seller.TID, "")
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessCancel(ctx, 0, "abc")
		if err == nil {
			t.Error(err)
			return
		}
	}
}

func TestFuturesMatcherBlewup(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		7: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on limit buy\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(96))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder1.OrderID)
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusCanceled)
		assertSmall()

		//will trigger to margin add on cancel
		buyOpenOrder3, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder3.OrderID)
		buyOpenOrder4, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(96))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder4.OrderID)
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder4.OrderID)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(buyOpenOrder4.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		//wlll trigger to margin free on place and margin add on cancel
		buyOpenOrder5, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder5.OrderID)
		assetOrderStatus(buyOpenOrder5.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
		assetOrderStatus(buyOpenOrder5.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		//will trigger to blow on sell
		buyOpenOrder6, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder6.OrderID)
		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on limit sell\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		}
		assertSmall()

		buyOpenOrder, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(98))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assertSmall()

		sellOpenOrder1, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assertSmall()
		sellOpenOrder2, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(103))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder1.OrderID)
		assertSmall()
		matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder1.OrderID)
		matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder2.OrderID)
		assetOrderStatus(sellOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetOrderStatus(sellOpenOrder2.OrderID, gexdb.OrderStatusCanceled)
		assertSmall()

		//will trigger to margin add on cancel
		sellOpenOrder3, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder3.OrderID)
		sellOpenOrder4, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(106))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder4.OrderID)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusPending)
		matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder4.OrderID)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder4.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		//wlll trigger to margin free on place and margin add on cancel
		sellOpenOrder5, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder5.OrderID)
		assetOrderStatus(sellOpenOrder5.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder5.OrderID)
		assetOrderStatus(sellOpenOrder5.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		//will trigger to blow on sell
		sellOpenOrder6, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(120))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder6.OrderID)
		buyOpenOrder, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder3.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on market\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		sellOpenOrder, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on cancel full one\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		sellCancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder1.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell cacnel order %v\n", sellCancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on cancel full one part\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		buyOpenOrder3, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder3.OrderID)
		sellCancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder1.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell cacnel order %v\n", sellCancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on cancel part one full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		sellCancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder1.OrderID)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell cacnel order %v\n", sellCancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup skip for missing market\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(99))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assertSmall()

		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
	}
}

func TestFuturesMatcherChangeLever(t *testing.T) {
	pgx.MockerStart()
	defer pgx.MockerStop()
	clear()
	enabled := map[int]bool{
		0: true,
		2: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: change lever\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		//
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		//
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		//
		err = matcher.ChangeLever(ctx, env.Buyer.TID, 5)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		//
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(10))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		err = matcher.ChangeLever(ctx, env.Buyer.TID, 5)
		if !IsErrOrderPending(err) {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("shoud pending err ->\n%v\n", ErrStack(err))
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder2.OrderID)
		//
		pgx.MockerClear()
		pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1, 2, "Tx.Exec", 1, 2, 3).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = matcher.ChangeLever(ctx, env.Buyer.TID, 1)
			return
		})
		pgx.ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			err = matcher.ChangeLever(ctx, env.Buyer.TID, 0)
			return
		})
	}
}

func TestFuturesMatcherError(t *testing.T) {
	clear()
	pgx.MockerStart()
	defer pgx.MockerStop()
	enabled := map[int]bool{
		0: true,
		4: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: market sell buy full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetAll := func() {
			assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
			assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
			assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerClear()

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder := &gexdb.Order{
			Type:       gexdb.OrderTypeTrigger,
			OrderID:    matcher.NewOrderID(),
			UserID:     env.Seller.TID,
			Symbol:     matcher.Symbol,
			Side:       gexdb.OrderSideBuy,
			TotalPrice: decimal.NewFromFloat(110),
			Status:     gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		pgx.MockerClear()
		pgx.MockerSetRangeCall("Rows.Scan", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
			return
		})
		_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: market sell part buy full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetAll := func() {
			assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(40.4))
			assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9959.6))
			assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		buyOpenOrder, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(110), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: limit sell buy full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetAll := func() {
			assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
			assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
			assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder := &gexdb.Order{
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			UserID:   env.Seller.TID,
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideBuy,
			Quantity: decimal.NewFromFloat(1),
			Price:    decimal.NewFromFloat(110),
			Status:   gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		pgx.MockerClear()
		pgx.MockerSetRangeCall("Rows.Scan", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
			return
		})
		_, err = matcher.ProcessOrder(ctx, sellCloseOrder)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: limit sell part buy full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetAll := func() {
			assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(40.4))
			assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9959.6))
			assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartialled)

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))

		fmt.Printf("%v start buy close order\n", env.Seller.TID)
		buyCloseOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(110))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v buy close order %v\n", env.Seller.TID, buyCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.22))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		fmt.Printf("%v start sell close order\n", env.Seller.TID)
		sellCloseOrder, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(110), decimal.Zero)
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("%v sell close order %v\n", env.Seller.TID, sellCloseOrder.OrderID)
		assetOrderStatus(buyCloseOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellCloseOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10009.58))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.58))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: limit sell buy full\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetAll := func() {
			assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20.2))
			assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
			assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1, "Tx.Query", 1, "Tx.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
			if err == nil {
				return
			}
			assetAll()
			return
		})

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Seller.TID, sellOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.Zero)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: cancel open\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		assetAll := func() {
			assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusPending)
			assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.2))
			assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
			assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9989.8))
			assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
		}
		assetAll()
		pgx.MockerClear()

		pgx.MockerSetCall("Pool.Begin", 1, "Tx.Commit", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerSetCall("Tx.Query", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
			if err == nil {
				return
			}
			// fmt.Println(ErrStack(err))
			assetAll()
			return
		})

		pgx.MockerSetRangeCall("Rows.Scan", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		pgx.MockerSetRangeCall("Tx.Exec", 1, 3).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
			if err == nil {
				return
			}
			assetAll()
			return
		})

		cancelOrder, err := matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(buyOpenOrder.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceLocked(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceMargin(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Buyer.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10000))
		assetHoldingAmount(env.Buyer.TID, futuresHoldingSymbol, decimal.Zero)
	}
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: order id error\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.NewOrderID = func() string {
			return "1"
		}
		buyOpenOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder.OrderID)
		_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on limit buy\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(96))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder1.OrderID)
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder2.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusCanceled)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusCanceled)
		assertSmall()

		//will trigger to margin add on cancel
		buyOpenOrder3, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder3.OrderID)
		buyOpenOrder4, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(96))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder4.OrderID)
		pgx.MockerClear()
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 10).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder4.OrderID)
			if err == nil {
				return
			}
			assertSmall()
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 10).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder4.OrderID)
			if err == nil {
				return
			}
			assertSmall()
			return
		})
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder4.OrderID)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(buyOpenOrder4.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		//wlll trigger to margin free on place and margin add on cancel
		pgx.MockerClear()
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
			return
		})
		pgx.MockerPanicRangeCall("Tx.Exec", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
			return
		})
		buyOpenOrder5, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(92))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder5.OrderID)
		assetOrderStatus(buyOpenOrder5.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		pgx.MockerClear()
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 10).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 10).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
			return
		})
		pgx.MockerPanicRangeCall("Tx.Exec", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
			return
		})
		matcher.ProcessCancel(ctx, env.Buyer.TID, buyOpenOrder5.OrderID)
		assetOrderStatus(buyOpenOrder5.OrderID, gexdb.OrderStatusCanceled)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(12))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(8.8))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))

		//will trigger to blow on sell
		buyOpenOrder6, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder6.OrderID)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusPending)
		pgx.MockerClear()
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 12).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 12).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
			return
		})
		pgx.MockerPanicRangeCall("Tx.Exec", 1, 8).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
			return
		})
		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup on market\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		buyOpenOrder3, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(50))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder3.OrderID)
		pgx.MockerClear()
		pgx.MockerPanicRangeCall("Tx.Query", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
			return
		})
		pgx.MockerSetRangeCall("Tx.Query", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
			return
		})
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 17).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 17).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
			return
		})
		sellOpenOrder, err = matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder3.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
	if testCount++; enabled[0] || enabled[testCount] {
		clear()
		fmt.Printf("\n\n==>start case %v: blowup order id error\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		matcher.MarginMax = decimal.NewFromFloat(0.9)
		matcher.MarginAdd = decimal.NewFromFloat(0.1)
		sellOpenOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		smallOpenOrder, err := matcher.ProcessLimit(ctx, env.Small.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("small open order %v\n", smallOpenOrder.OrderID)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(smallOpenOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceMargin(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceLocked(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(20))
		assetBalanceFree(env.Seller.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(9979.8))
		assetHoldingAmount(env.Seller.TID, futuresHoldingSymbol, decimal.NewFromFloat(-1))
		assertSmall := func() {
			assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10))
			assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(10.8))
			assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(1))
		}
		assertSmall()

		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(102))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assertSmall()

		buyOpenOrder1, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder1.OrderID)
		assertSmall()

		//will trigger to blow on sell
		buyOpenOrder2, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(60))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("buy open order %v\n", buyOpenOrder2.OrderID)
		matcher.NewOrderID = func() string { return sellOpenOrder.OrderID }
		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err == nil {
			t.Error(ErrStack(err))
			return
		}
		matcher.NewOrderID = gexdb.NewOrderID
		sellOpenOrder, err = matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(95))
		if err != nil {
			t.Error(ErrStack(err))
			return
		}
		fmt.Printf("sell open order %v\n", sellOpenOrder.OrderID)
		assetOrderStatus(buyOpenOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOpenOrder2.OrderID, gexdb.OrderStatusPending)
		assetOrderStatus(sellOpenOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceMargin(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceFree(env.Small.TID, env.Area, futuresBalanceQuote, decimal.NewFromFloat(0))
		assetHoldingAmount(env.Small.TID, futuresHoldingSymbol, decimal.NewFromFloat(0))
	}
}

func TestFuturesMatcherDepth(t *testing.T) {
	clear()
	enabled := map[int]bool{
		0: true,
		6: true,
	}
	testCount := 0
	if testCount++; enabled[0] || enabled[testCount] {
		fmt.Printf("\n\n==>start case %v: depth cancel\n", testCount)
		//
		env := testFuturesInit(testCount)
		matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
		for i := 0; i < 10; i++ {
			buyOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(10-int64(i))), decimal.NewFromFloat(99.5).Sub(decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(int64(i)+1))))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Printf("buy order %v\n", buyOrder.OrderID)
			assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		}
		for i := 0; i < 10; i++ {
			sellOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(10-int64(i))), decimal.NewFromFloat(100.5).Add(decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(int64(i)+1))))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Printf("sell order %v\n", sellOrder.OrderID)
			assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		}
		assetDepthMust(matcher.Depth(0), 10, 10)

		buyOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(99.6))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetDepthMust(matcher.Depth(0), 11, 10)

		_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		assetDepthMust(matcher.Depth(0), 10, 10)
		// assetBalanceLocked(userQuote.TID, spotBalanceQuote, decimal.NewFromFloat(100))
	}
}

func TestFuturesMatcherParallel(t *testing.T) {
	// clear()
	env := testFuturesInit(0)
	matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
	matcher.PrecisionPrice = 8
	matcher.PrecisionQuantity = 8
	// for _, i := range []int64{1, 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20} {
	elapsed, avg := ParallelTest(100, 5, func(i int64) {
		switch i % 9 {
		case 0:
			sellOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(5).Mul(decimal.NewFromFloat(0.01))))
			if err != nil {
				panic(ErrStack(err))
			}
			buyOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(1).Mul(decimal.NewFromFloat(0.01))))
			if err != nil {
				panic(ErrStack(err))
			}
			_, err = matcher.ProcessCancel(ctx, env.Seller.TID, sellOrder.OrderID)
			if err != nil && !IsErrNotCancelable(err) {
				panic(ErrStack(err))
			}
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOrder.OrderID)
			if err != nil && !IsErrNotCancelable(err) {
				panic(ErrStack(err))
			}
		case 1, 2, 3:
			_, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(4+i%3)).Mul(decimal.NewFromInt(1).Mul(decimal.NewFromFloat(0.01)))))
			if err != nil {
				panic(ErrStack(err))
			}
		case 4, 5, 6:
			_, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(3+i%3)).Mul(decimal.NewFromInt(1).Mul(decimal.NewFromFloat(0.01)))))
			if err != nil {
				panic(ErrStack(err))
			}
		case 7:
			_, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.01))
			if err != nil {
				panic(ErrStack(err))
			}
		case 8:
			_, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(0.01))
			if err != nil {
				panic(ErrStack(err))
			}
		}
	})
	fmt.Printf("\n\nElapsed: %s\nTransactions per second (avg): %f\n", elapsed, avg)
}

func BenchmarkFuturesMatcher(b *testing.B) {
	clear()
	env := testFuturesInit(0)
	matcher := NewFuturesMatcher(futuresHoldingSymbol, futuresBalanceQuote, env.Monitor)
	matcher.PrecisionPrice = 8
	matcher.PrecisionQuantity = 8
	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		switch i % 9 {
		case 0:
			sellOrder, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(5)))
			if err != nil {
				panic(err)
			}
			buyOrder, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(1)))
			if err != nil {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, env.Seller.TID, sellOrder.OrderID)
			if err != nil {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, env.Buyer.TID, buyOrder.OrderID)
			if err != nil {
				panic(err)
			}
		case 1, 2, 3:
			_, err := matcher.ProcessLimit(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(4+i%3))))
			if err != nil {
				panic(ErrStack(err))
			}
		case 4, 5, 6:
			_, err := matcher.ProcessLimit(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(3+i%3))))
			if err != nil {
				panic(ErrStack(err))
			}
		case 7:
			_, err := matcher.ProcessMarket(ctx, env.Buyer.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.01))
			if err != nil {
				panic(ErrStack(err))
			}
		case 8:
			_, err := matcher.ProcessMarket(ctx, env.Seller.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(0.01))
			if err != nil {
				panic(ErrStack(err))
			}
		}
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second (avg): %f\n", elapsed, float64(b.N*32)/elapsed.Seconds())
}
