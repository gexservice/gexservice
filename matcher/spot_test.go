package matcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

const (
	spotBalanceBase   = "YWE"
	spotBalanceQuote  = "USDT"
	spotBalanceSymbol = "spot.YWEUSDT"
)

var spotFeeRate = decimal.NewFromFloat(0.002)
var spotBalanceAll = []string{spotBalanceBase, spotBalanceQuote}

func TestSpotMatcherBootstrap(t *testing.T) {
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherLimit-Base")
	userQuote := testAddUser("TestSpotMatcherLimit-Quote")
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
	startBalances, err := gexdb.CountBalance(ctx, area, time.Time{}, time.Now())
	if err != nil || !startBalances[spotBalanceBase].Equal(decimal.NewFromFloat(1000)) || !startBalances[spotBalanceQuote].Equal(decimal.NewFromFloat(1000)) {
		t.Error(err)
		return
	}
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	changed, err := matcher.Bootstrap(ctx)
	if err != nil || len(changed.Orders) > 0 {
		t.Error(err)
		return
	}

	buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("buy order %v\n", buyOrder.OrderID)
	assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
	assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))
	sellOrder1, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("sell order %v\n", sellOrder1.OrderID)
	sellOrder2, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(110))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("sell order %v\n", sellOrder2.OrderID)
	assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
	assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusDone)
	assetOrderStatus(sellOrder2.OrderID, gexdb.OrderStatusPending)

	pgx.MockerStart()
	defer pgx.MockerStop()

	matcher = NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	pgx.MockerPanicCall("Pool.Begin", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = matcher.Bootstrap(ctx)
		return
	})
	pgx.MockerSetCall("Pool.Begin", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = matcher.Bootstrap(ctx)
		return
	})
	pgx.MockerSetRangeCall("Tx.Exec", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
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

func TestSpotMatcherMarket(t *testing.T) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherMarket-Base")
	userQuote := testAddUser("TestSpotMatcherMarket-Quote")
	userNone := testAddUser("TestSpotMatcherMarket-NONE")
	_, err := gexdb.TouchBalance(ctx, area, spotBalanceAll, userBase.TID, userQuote.TID, userNone.TID)
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
	startBalances, err := gexdb.CountBalance(ctx, area, time.Time{}, time.Now())
	if err != nil || !startBalances[spotBalanceBase].Equal(decimal.NewFromFloat(10000)) || !startBalances[spotBalanceQuote].Equal(decimal.NewFromFloat(10000)) {
		t.Error(err)
		return
	}
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	{ //sell buy all, invest
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0.5))

		buyOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(50), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceFree(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(9950))
		assetBalanceFree(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0.5).Mul(decimal.NewFromFloat(1).Sub(spotFeeRate)))
		assetBalanceFree(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(50).Mul(decimal.NewFromFloat(1).Sub(spotFeeRate)))
		assetBalanceFree(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(9999.5))
	}
	{ //sell buy all, invest
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0.5))

		buyOrder := &gexdb.Order{
			UserID:     userQuote.TID,
			Type:       gexdb.OrderTypeTrigger,
			OrderID:    matcher.NewOrderID(),
			Symbol:     matcher.Symbol,
			Side:       gexdb.OrderSideBuy,
			TotalPrice: decimal.NewFromFloat(50),
			Status:     gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, buyOrder)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessOrder(ctx, buyOrder)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceFree(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(9900))
		assetBalanceFree(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(1).Mul(decimal.NewFromFloat(1).Sub(spotFeeRate)))
		assetBalanceFree(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(100).Mul(decimal.NewFromFloat(1).Sub(spotFeeRate)))
		assetBalanceFree(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(9999))
	}
	{ //sell buy all, quantity
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //sell buy partial, invest
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder1, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 1 %v\n", buyOrder1.OrderID)
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPartialled)

		buyOrder2, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 2 %v\n", buyOrder2.OrderID)
		assetOrderStatus(buyOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //sell buy partial, quantity
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder1, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 1 %v\n", buyOrder1.OrderID)
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPartialled)

		buyOrder2, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 2 %v\n", buyOrder2.OrderID)
		assetOrderStatus(buyOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //buy sell all
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		sellOrder, err := matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //buy sell partial, buy partial
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))

		sellOrder1, err := matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 1 %v\n", sellOrder1.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusDone)

		sellOrder2, err := matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 2 %v\n", sellOrder2.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //buy sell partial, sell partial
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		sellOrder1, err := matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(2))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 1 %v\n", sellOrder1.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	//test error
	{ //buy not found
		buyOrder1, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusCanceled)

		buyOrder2, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		assetOrderStatus(buyOrder2.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //sell not found
		sellOrder, err := matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //buy not enought, invest
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if !IsErrBalanceNotEnought(err) {
			t.Error(ErrStack(err))
			return
		}
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //buy not enought, quantity
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if !IsErrBalanceNotEnought(err) {
			t.Error(err)
			return
		}
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //sell not enought
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if !IsErrBalanceNotEnought(err) {
			t.Error(err)
			return
		}
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		sellOrder, err := matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //arg error
		_, err = matcher.ProcessMarket(ctx, 0, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, userNone.TID, "0", decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideSell, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}

		//
		doneOrder := &gexdb.Order{
			UserID:   userBase.TID,
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideSell,
			Quantity: decimal.NewFromFloat(0.5),
			Status:   gexdb.OrderStatusDone,
		}
		err = gexdb.AddOrder(ctx, doneOrder)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessOrder(ctx, doneOrder)
		if err == nil {
			t.Error(err)
			return
		}

		matcher.PrepareProcess = func(ctx context.Context, matcher *SpotMatcher, userID int64) error { return fmt.Errorf("xxxx") }
		_, err = matcher.ProcessMarket(ctx, userNone.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
	}
}

func TestSpotMatcherLimit(t *testing.T) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherLimit-Base")
	userQuote := testAddUser("TestSpotMatcherLimit-Quote")
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
	startBalances, err := gexdb.CountBalance(ctx, area, time.Time{}, time.Now())
	if err != nil || !startBalances[spotBalanceBase].Equal(decimal.NewFromFloat(1000)) || !startBalances[spotBalanceQuote].Equal(decimal.NewFromFloat(1000)) {
		t.Error(err)
		return
	}
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	{ //buy sell all
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(50))

		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //buy sell all, prepare
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(50))

		sellOrder := &gexdb.Order{
			UserID:   userBase.TID,
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideSell,
			Quantity: decimal.NewFromFloat(0.5),
			Price:    decimal.NewFromFloat(100),
			Status:   gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, sellOrder)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessOrder(ctx, sellOrder)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell buy all
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //buy sell, buy partial
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))

		sellOrder1, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 1 %v\n", sellOrder1.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusDone)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))

		sellOrder2, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 2 %v\n", sellOrder2.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //buy sell, sell partial
		buyOrder1, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 1 %v\n", buyOrder1.OrderID)
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		sellOrder1, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 1 %v\n", sellOrder1.OrderID)
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusPartialled)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder2, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 2 %v\n", buyOrder2.OrderID)
		assetOrderStatus(buyOrder2.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell buy, buy partial
		sellOrder1, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 1 %v\n", sellOrder1.OrderID)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(sellOrder1.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))

		sellOrder2, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order 2 %v\n", sellOrder2.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell buy, sell partial
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder1, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 1 %v\n", buyOrder1.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(buyOrder1.OrderID, gexdb.OrderStatusDone)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder2, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order 2 %v\n", buyOrder2.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOrder2.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	doneBalances, err := gexdb.CountBalance(ctx, gexdb.BalanceAreaSpot, time.Time{}, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	doneFee, err := gexdb.CountOrderFee(ctx, time.Time{}, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	if !startBalances[spotBalanceBase].Equal(doneBalances[spotBalanceBase].Add(doneFee[spotBalanceBase])) {
		t.Error("Base balance fail")
		return
	}
	if !startBalances[spotBalanceQuote].Equal(doneBalances[spotBalanceQuote].Add(doneFee[spotBalanceQuote])) {
		t.Error("Base balance fail")
		return
	}
	//test error
	{ //buy sell too much
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(10000), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
	}
	{ //arg error
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, "0", decimal.Zero, decimal.NewFromFloat(1))
		if err == nil {
			t.Error(err)
			return
		}

		//
		doneOrder := &gexdb.Order{
			UserID:   userBase.TID,
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideSell,
			Quantity: decimal.NewFromFloat(0.5),
			Price:    decimal.NewFromFloat(100),
			Status:   gexdb.OrderStatusDone,
		}
		err = gexdb.AddOrder(ctx, doneOrder)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessOrder(ctx, doneOrder)
		if err == nil {
			t.Error(err)
			return
		}

		//prepare erro
		matcher.PrepareProcess = func(ctx context.Context, matcher *SpotMatcher, userID int64) error { return fmt.Errorf("xxxx") }
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
	}
}

func TestSpotMatcherCancel(t *testing.T) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherCancel-Base")
	userQuote := testAddUser("TestSpotMatcherCancel-Quote")
	_, err := gexdb.TouchBalance(ctx, area, spotBalanceAll, userBase.TID, userQuote.TID)
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
	startBalances, err := gexdb.CountBalance(ctx, area, time.Time{}, time.Now())
	if err != nil || !startBalances[spotBalanceBase].Equal(decimal.NewFromFloat(10000)) || !startBalances[spotBalanceQuote].Equal(decimal.NewFromFloat(10000)) {
		t.Error(err)
		return
	}
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	{ //buy cancel
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		cancelOrder, err := matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)
		cancelOrder, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if !IsErrNotCancelable(err) {
			t.Error(err)
			return
		}
		fmt.Println("-->", ErrStack(err))
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell cancel
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		cancelOrder, err := matcher.ProcessCancel(ctx, sellOrder.UserID, sellOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //buy sell, buy cancel partial
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))

		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))

		cancelOrder, err := matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //buy sell, sell cancel partial
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPartialled)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		cancelOrder, err := matcher.ProcessCancel(ctx, sellOrder.UserID, sellOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell buy, sell cancel partial
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))

		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPartialled)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(2))
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))

		cancelOrder, err := matcher.ProcessCancel(ctx, sellOrder.UserID, sellOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //sell buy, buy cancel partial
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))

		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPartialled)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(200))

		cancelOrder, err := matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusPartCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //cancel not access
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))

		_, err = matcher.ProcessCancel(ctx, 10, buyOrder.OrderID)
		if err == nil {
			t.Error(err)
			return
		}
		cancelOrder, err := matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)
		cancelOrder, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if !IsErrNotCancelable(err) {
			t.Error(err)
			return
		}
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetDepthEmpty(matcher.Depth(0))
	}
	{ //multi, buy cancel
		for i := 0; i < 10; i++ {
			buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(10-int64(i))), decimal.NewFromFloat(99.5).Sub(decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(int64(i)+1))))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Printf("buy order %v\n", buyOrder.OrderID)
			assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		}
		for i := 0; i < 10; i++ {
			sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(10-int64(i))), decimal.NewFromFloat(100.5).Add(decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(int64(i)+1))))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Printf("sell order %v\n", sellOrder.OrderID)
			assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		}
		assetDepthMust(matcher.Depth(0), 10, 10)

		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.5), decimal.NewFromFloat(99.6))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetDepthMust(matcher.Depth(0), 11, 10)

		_, err = matcher.ProcessCancel(ctx, userQuote.TID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		assetDepthMust(matcher.Depth(0), 10, 10)
		// assetBalanceLocked(userQuote.TID, spotBalanceQuote, decimal.NewFromFloat(100))
	}
	{ //arg error
		_, err := matcher.ProcessCancel(ctx, userBase.TID, "")
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
	doneBalances, err := gexdb.CountBalance(ctx, gexdb.BalanceAreaSpot, time.Time{}, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	doneFee, err := gexdb.CountOrderFee(ctx, time.Time{}, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	if !startBalances[spotBalanceBase].Equal(doneBalances[spotBalanceBase].Add(doneFee[spotBalanceBase])) {
		t.Error("Base balance fail")
		return
	}
	if !startBalances[spotBalanceQuote].Equal(doneBalances[spotBalanceQuote].Add(doneFee[spotBalanceQuote])) {
		t.Error("Base balance fail")
		return
	}
}

func TestSpotMatcherError(t *testing.T) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherMarketError-Base")
	userQuote := testAddUser("TestSpotMatcherMarketError-Quote")
	// userNone := testAddUser("TestSpotMatcher-NONE")
	_, err := gexdb.TouchBalance(ctx, gexdb.BalanceAreaSpot, spotBalanceAll, userBase.TID, userQuote.TID)
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
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	pgx.MockerStart()
	defer pgx.MockerStop()
	{ //market buy
		sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("sell order %v\n", sellOrder.OrderID)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(1))
		pgx.MockerClear()

		pgx.MockerSet("Pool.Begin", 1)
		_, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err == nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 5).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			return
		})
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 7).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
			return
		})

		buyOrder := &gexdb.Order{
			UserID:     userQuote.TID,
			Type:       gexdb.OrderTypeTrigger,
			OrderID:    matcher.NewOrderID(),
			Symbol:     matcher.Symbol,
			Side:       gexdb.OrderSideBuy,
			TotalPrice: decimal.NewFromFloat(100),
			Status:     gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, buyOrder)
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()
		pgx.MockerSetRangeCall("Rows.Scan", 1, 3).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessOrder(ctx, buyOrder)
			return
		})

		buyOrder, err = matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(100), decimal.Zero)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)
		assetOrderStatus(sellOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		pgx.MockerClear()
	}
	{ //limit buy
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(2), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()

		pgx.MockerSet("Pool.Begin", 1)
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2.5), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 6).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2.5), decimal.NewFromFloat(100))
			assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(3))
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 6).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2.5), decimal.NewFromFloat(100))
			assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(3))
			return
		})
		pgx.MockerPanicRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2.5), decimal.NewFromFloat(100))
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 9).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(2.5), decimal.NewFromFloat(100))
			return
		})

		buyOrder := &gexdb.Order{
			UserID:   userQuote.TID,
			Type:     gexdb.OrderTypeTrigger,
			OrderID:  matcher.NewOrderID(),
			Symbol:   matcher.Symbol,
			Side:     gexdb.OrderSideBuy,
			Quantity: decimal.NewFromFloat(3),
			Price:    decimal.NewFromFloat(100),
			Status:   gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, buyOrder)
		if err != nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()
		pgx.MockerSetRangeCall("Rows.Scan", 1, 3).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessOrder(ctx, buyOrder)
			return
		})

		buyOrder, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(3), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusDone)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		pgx.MockerClear()
	}
	{ //cancel error
		buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("buy order %v\n", buyOrder.OrderID)
		assetOrderStatus(buyOrder.OrderID, gexdb.OrderStatusPending)
		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(100))
		pgx.MockerClear()

		pgx.MockerSet("Pool.Begin", 1)
		_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err == nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()

		pgx.MockerPanicRangeCall("Rows.Scan", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
			return
		})
		pgx.MockerSetRangeCall("Rows.Scan", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
			return
		})

		pgx.MockerPanicRangeCall("Tx.Exec", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
			return
		})
		pgx.MockerSetRangeCall("Tx.Exec", 1, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
			_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
			return
		})

		pgx.MockerSet("Tx.Commit", 1)
		_, err = matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err == nil {
			t.Error(err)
			return
		}
		pgx.MockerClear()

		cancelOrder, err := matcher.ProcessCancel(ctx, buyOrder.UserID, buyOrder.OrderID)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("cancel order %v\n", cancelOrder.OrderID)
		assetOrderStatus(cancelOrder.OrderID, gexdb.OrderStatusCanceled)

		assetBalanceLocked(userQuote.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userQuote.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceQuote, decimal.NewFromFloat(0))
		assetBalanceLocked(userBase.TID, area, spotBalanceBase, decimal.NewFromFloat(0))
	}
	{ //order id error
		matcher.NewOrderID = func() string {
			return "1"
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err == nil {
			t.Error(err)
			return
		}
		matcher.NewOrderID = gexdb.NewOrderID
		_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestSpotMatcherDepth(t *testing.T) {
	clear()
	userBase := testAddUser("TestSpotMatcherDepth-Base")
	userQuote := testAddUser("TestSpotMatcherDepth-Quote")
	// userNone := testAddUser("TestSpotMatcher-NONE")
	_, err := gexdb.TouchBalance(ctx, gexdb.BalanceAreaSpot, spotBalanceAll, userBase.TID, userQuote.TID)
	if err != nil {
		t.Error(err)
		return
	}
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userBase.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceBase,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
		UserID: userQuote.TID,
		Area:   gexdb.BalanceAreaSpot,
		Asset:  spotBalanceQuote,
		Free:   decimal.NewFromFloat(10000),
		Status: gexdb.BalanceStatusNormal,
	})
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	_, err = matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(1), decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1), decimal.NewFromFloat(90))
	if err != nil {
		t.Error(err)
		return
	}
	depth := matcher.Depth(0)
	if len(depth.Bids) != 1 && len(depth.Asks) != 1 {
		t.Error("error")
		return
	}
}

func TestSpotMatcherParallel(t *testing.T) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherMarket-Base")
	userQuote := testAddUser("TestSpotMatcherMarket-Quote")
	userNone := testAddUser("TestSpotMatcherMarket-NONE")
	_, err := gexdb.TouchBalance(ctx, area, spotBalanceAll, userBase.TID, userQuote.TID, userNone.TID)
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
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	matcher.Timeout = time.Hour
	elapsed, avg := ParallelTest(100, 5, func(i int64) {
		switch i % 9 {
		case 0:
			sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.001), decimal.NewFromFloat(1).Mul(decimal.NewFromInt(5)))
			if err != nil {
				panic(err)
			}
			buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.001), decimal.NewFromFloat(1).Mul(decimal.NewFromInt(1)))
			if err != nil {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, userBase.TID, sellOrder.OrderID)
			if err != nil && !IsErrNotCancelable(err) {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, userQuote.TID, buyOrder.OrderID)
			if err != nil && !IsErrNotCancelable(err) {
				panic(err)
			}
		case 1, 2, 3:
			matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.001), decimal.NewFromFloat(1).Mul(decimal.NewFromInt(int64(4+i%3))))
		case 4, 5, 6:
			matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.001), decimal.NewFromFloat(1).Mul(decimal.NewFromInt(int64(3+i%3))))
		case 7:
			matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.001))
		case 8:
			matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(0.001))
		}
	})
	fmt.Printf("\n\nElapsed: %s\nTransactions per second (avg): %f\n", elapsed, avg)
}

func BenchmarkSpotMatcher(b *testing.B) {
	clear()
	area := gexdb.BalanceAreaSpot
	userBase := testAddUser("TestSpotMatcherMarket-Base")
	userQuote := testAddUser("TestSpotMatcherMarket-Quote")
	userNone := testAddUser("TestSpotMatcherMarket-NONE")
	_, err := gexdb.TouchBalance(ctx, area, spotBalanceAll, userBase.TID, userQuote.TID, userNone.TID)
	if err != nil {
		b.Error(err)
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
	matcher := NewSpotMatcher(spotBalanceSymbol, spotBalanceBase, spotBalanceQuote, MatcherMonitorF(func(ctx context.Context, event *MatcherEvent) {
	}))
	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		switch i % 9 {
		case 0:
			sellOrder, err := matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(5)))
			if err != nil {
				panic(err)
			}
			buyOrder, err := matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(1)))
			if err != nil {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, userBase.TID, sellOrder.OrderID)
			if err != nil {
				panic(err)
			}
			_, err = matcher.ProcessCancel(ctx, userQuote.TID, buyOrder.OrderID)
			if err != nil {
				panic(err)
			}
		case 1, 2, 3:
			matcher.ProcessLimit(ctx, userBase.TID, gexdb.OrderSideSell, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(4+i%3))))
		case 4, 5, 6:
			matcher.ProcessLimit(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(0.01), decimal.NewFromFloat(100).Mul(decimal.NewFromInt(int64(3+i%3))))
		case 7:
			matcher.ProcessMarket(ctx, userQuote.TID, gexdb.OrderSideBuy, decimal.Zero, decimal.NewFromFloat(0.01))
		case 8:
			matcher.ProcessMarket(ctx, userBase.TID, gexdb.OrderSideSell, decimal.Zero, decimal.NewFromFloat(0.01))
		}
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second (avg): %f\n", elapsed, float64(b.N*32)/elapsed.Seconds())
}
