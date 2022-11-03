package gexdb

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/codingeasygo/util/xsql"
// 	"github.com/shopspring/decimal"
// )

// func TestTopupOrder(t *testing.T) {
// 	user := testAddUser("TestTopupOrder")
// 	TouchBalance(BalanceAssetAll, user.TID)
// 	order := &Order{
// 		Type:      OrderTypeTopup,
// 		UserID:    user.TID,
// 		InBalance: BalanceAssetMMK,
// 		InFilled:  decimal.NewFromFloat(100),
// 		Status:    OrderStatusPending,
// 	}
// 	err := CreateOrder(order)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = UpdateOrderPrepay(order.TID, xsql.M{"test": 1})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = PayTopupOrder(order.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	balances, err := ListUserBalance(order.UserID, BalanceAssetAll, BalanceStatusAll)
// 	if err != nil || len(balances) != 2 {
// 		t.Error(err)
// 		return
// 	}
// 	if !balances[BalanceAssetMMK].Free.Equal(order.InFilled) {
// 		t.Error("value error")
// 		return
// 	}

// 	findOrder, err := FindOrderByOrderID(order.OrderID)
// 	if err != nil || findOrder.TID != order.TID {
// 		t.Error(err)
// 		return
// 	}
// 	findOrder, err = FindOrderByOrderID(fmt.Sprintf("%v", order.TID))
// 	if err != nil || findOrder.TID != order.TID {
// 		t.Error(err)
// 		return
// 	}

// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	order2 := &Order{
// 		Type:      OrderTypeTopup,
// 		UserID:    user.TID,
// 		InBalance: BalanceAssetMMK,
// 		InFilled:  decimal.NewFromFloat(100),
// 		Status:    OrderStatusPending,
// 	}
// 	err = CreateOrder(order2)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	order3 := &Order{
// 		Type:      OrderTypeGoldbar,
// 		UserID:    user.TID,
// 		InBalance: BalanceAssetMMK,
// 		InFilled:  decimal.NewFromFloat(100),
// 		Status:    OrderStatusPending,
// 	}
// 	err = CreateOrder(order3)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//
// 	//PayTopupOrder error
// 	pgx.MockerSet("Pool.Begin", 1)
// 	err = PayTopupOrder(order2.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1})
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 1)
// 	err = PayTopupOrder(order2.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1})
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Tx.Exec", 1)
// 	err = PayTopupOrder(order2.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1})
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 2)
// 	err = PayTopupOrder(order2.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1})
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	err = PayTopupOrder(order.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1}) //status error
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	err = PayTopupOrder(order3.OrderID, decimal.NewFromFloat(100), xmap.M{"test": 1}) //type error
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// }
