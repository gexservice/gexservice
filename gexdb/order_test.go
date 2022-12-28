package gexdb

import (
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

func TestOrder(t *testing.T) {
	clear()
	feeAll, _, err := CountOrderFee(ctx, OrderAreaSpot, time.Time{}, time.Now())
	if err != nil || len(feeAll) > 0 {
		t.Error(err)
		return
	}
	volumeAll, err := CountOrderVolume(ctx, OrderAreaSpot, time.Time{}, time.Now())
	if err != nil || len(volumeAll) > 0 {
		t.Error(err)
		return
	}
	user := testAddUser("TestOrder")
	order := &Order{
		Type:       OrderTypeTrade,
		UserID:     user.TID,
		Creator:    user.TID,
		Area:       OrderAreaSpot,
		Symbol:     "test",
		OrderID:    NewOrderID(),
		FeeBalance: "test",
		FeeFilled:  decimal.NewFromFloat(1),
		Status:     OrderStatusDone,
	}
	err = AddOrder(ctx, order)
	if err != nil {
		t.Error(err)
		return
	}
	findOrder, err := FindOrderByOrderID(ctx, 0, order.OrderID)
	if err != nil || order.TID != findOrder.TID {
		t.Error(err)
		return
	}
	findOrder, err = FindOrderByOrderIDCall(Pool(), ctx, 0, order.OrderID, true)
	if err != nil || order.TID != findOrder.TID {
		t.Error(err)
		return
	}
	searcher := &OrderUnifySearcher{}
	searcher.Where.UserID = xsql.Int64Array{order.UserID}
	searcher.Where.Area = OrderAreaAll
	searcher.Where.Symbol = order.Symbol
	searcher.Where.Key = order.OrderID
	err = searcher.Apply(ctx)
	if err != nil || searcher.Count.Total != 1 || len(searcher.Query.Orders) != 1 {
		t.Errorf("%v,%v", err, searcher.Count.Total)
		return
	}
	_, err = ClearCanceledOrder(ctx, order.UserID, order.Symbol, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	having, err := CountPendingOrderCall(Pool(), ctx, order.UserID, order.Symbol)
	if err != nil || having > 0 {
		t.Error(err)
		return
	}
	orders, err := ListPendingOrder(ctx, order.UserID, "te", order.Symbol)
	if err != nil || len(orders) > 0 {
		t.Error(err)
		return
	}
	feeAll, _, err = CountOrderFee(ctx, order.Area, time.Time{}, time.Now())
	if err != nil || len(feeAll) < 1 {
		t.Error(err)
		return
	}
	volumeAll, err = CountOrderVolume(ctx, order.Area, time.Time{}, time.Now())
	if err != nil || len(volumeAll) < 1 {
		t.Error(err)
		return
	}
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()

}

func TestTriggerOrder(t *testing.T) {
	clear()
	user := testAddUser("TestTriggerOrder")
	symbol := "test"

	err := AddOrder(ctx, &Order{
		Symbol:       symbol,
		Type:         OrderTypeTrigger,
		UserID:       user.TID,
		Creator:      user.TID,
		OrderID:      NewOrderID(),
		Side:         OrderSideSell,
		Quantity:     decimal.NewFromFloat(1),
		Price:        decimal.NewFromFloat(100),
		TriggerType:  OrderTriggerTypeStopProfit,
		TriggerPrice: decimal.NewFromFloat(100),
		Status:       OrderStatusWaiting,
	})
	if err != nil {
		t.Error(err)
		return
	}
	searcher := &OrderUnifySearcher{}
	searcher.Where.Area = OrderAreaAll
	searcher.Where.TriggerType = OrderTriggerTypeArray{OrderTriggerTypeStopProfit}
	err = searcher.Apply(ctx)
	if err != nil || searcher.Count.Total != 1 || len(searcher.Query.Orders) != 1 {
		t.Errorf("%v,%v", err, searcher.Count.Total)
		return
	}
	orders, err := ListOrderForTrigger(ctx, symbol, decimal.Zero, decimal.NewFromFloat(100))
	if err != nil || len(orders) != 1 {
		t.Errorf("%v,%v", err, converter.JSON(orders))
		return
	}
	if orders[0].Side != OrderSideSell || orders[0].TriggerType != OrderTriggerTypeStopProfit {
		t.Error("data")
		return
	}
	updated, err := CancelTriggerOrder(ctx, user.TID, symbol, orders[0].TID)
	if err != nil || updated < 1 {
		t.Error(err)
		return
	}

	err = AddOrder(ctx, &Order{
		Symbol:       symbol,
		Type:         OrderTypeTrigger,
		UserID:       user.TID,
		Creator:      user.TID,
		OrderID:      NewOrderID(),
		Side:         OrderSideSell,
		Quantity:     decimal.NewFromFloat(1),
		Price:        decimal.NewFromFloat(80),
		TriggerType:  OrderTriggerTypeStopLoss,
		TriggerPrice: decimal.NewFromFloat(80),
		Status:       OrderStatusWaiting,
	})
	if err != nil {
		t.Error(err)
		return
	}
	orders, err = ListOrderForTrigger(ctx, symbol, decimal.Zero, decimal.NewFromFloat(80))
	if err != nil || len(orders) != 1 {
		t.Errorf("%v,%v", err, converter.JSON(orders))
		return
	}
	if orders[0].Side != OrderSideSell || orders[0].TriggerType != OrderTriggerTypeStopLoss {
		t.Error("data")
		return
	}
	updated, err = CancelTriggerOrder(ctx, user.TID, symbol, orders[0].TID)
	if err != nil || updated < 1 {
		t.Error(err)
		return
	}

	err = AddOrder(ctx, &Order{
		Symbol:       symbol,
		Type:         OrderTypeTrigger,
		UserID:       user.TID,
		Creator:      user.TID,
		OrderID:      NewOrderID(),
		Side:         OrderSideBuy,
		Quantity:     decimal.NewFromFloat(1),
		Price:        decimal.NewFromFloat(80),
		TriggerType:  OrderTriggerTypeStopProfit,
		TriggerPrice: decimal.NewFromFloat(80),
		Status:       OrderStatusWaiting,
	})
	if err != nil {
		t.Error(err)
		return
	}
	orders, err = ListOrderForTrigger(ctx, symbol, decimal.NewFromFloat(80), decimal.Zero)
	if err != nil || len(orders) != 1 {
		t.Errorf("%v,%v", err, converter.JSON(orders))
		return
	}
	if orders[0].Side != OrderSideBuy || orders[0].TriggerType != OrderTriggerTypeStopProfit {
		t.Error("data")
		return
	}
	updated, err = CancelTriggerOrder(ctx, user.TID, symbol, orders[0].TID)
	if err != nil || updated < 1 {
		t.Error(err)
		return
	}

	err = AddOrder(ctx, &Order{
		Symbol:       symbol,
		Type:         OrderTypeTrigger,
		UserID:       user.TID,
		Creator:      user.TID,
		OrderID:      NewOrderID(),
		Side:         OrderSideBuy,
		Quantity:     decimal.NewFromFloat(1),
		Price:        decimal.NewFromFloat(100),
		TriggerType:  OrderTriggerTypeStopLoss,
		TriggerPrice: decimal.NewFromFloat(100),
		Status:       OrderStatusWaiting,
	})
	if err != nil {
		t.Error(err)
		return
	}
	orders, err = ListOrderForTrigger(ctx, symbol, decimal.NewFromFloat(100), decimal.Zero)
	if err != nil || len(orders) != 1 {
		t.Errorf("%v,%v", err, converter.JSON(orders))
		return
	}
	if orders[0].Side != OrderSideBuy || orders[0].TriggerType != OrderTriggerTypeStopLoss {
		t.Error("data")
		return
	}
	updated, err = CancelTriggerOrder(ctx, user.TID, symbol, orders[0].TID)
	if err != nil || updated < 1 {
		t.Error(err)
		return
	}
}
