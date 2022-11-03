package gexdb

// import (
// 	"testing"
// 	"time"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/shopspring/decimal"
// )

// func TestOrderComm(t *testing.T) {
// 	order := &Order{
// 		UserID: 10,
// 		Status: OrderStatusDone,
// 	}
// 	err := CreateOrder(order)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	comms := []*OrderComm{
// 		{
// 			OrderID:   order.TID,
// 			UserID:    0,
// 			Type:      OrderCommTypeSys,
// 			InBalance: BalanceAssetYWE,
// 			InFee:     decimal.NewFromFloat(20),
// 			Status:    OrderCommStatusNormal,
// 		},
// 		{
// 			OrderID:   order.TID,
// 			UserID:    100,
// 			Type:      OrderCommTypeBroker,
// 			InBalance: BalanceAssetYWE,
// 			InFee:     decimal.NewFromFloat(80),
// 			Status:    OrderCommStatusNormal,
// 		},
// 	}
// 	err = AddOrderCommCall(Pool(), comms...)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	listComms, _, err := ListUserOrderComm(100, order.TID)
// 	if err != nil || len(listComms) != 1 {
// 		t.Error("error")
// 		return
// 	}
// 	countComms, err := CountUserOrderComm(100, time.Now().Add(-10*time.Hour), time.Now())
// 	if err != nil || len(countComms) != 1 {
// 		t.Error(err)
// 		return
// 	}
// 	countUserComms, err := CountMyUserOrderComm(100, 10)
// 	if err != nil || len(countUserComms) != 1 {
// 		t.Error(err)
// 		return
// 	}
// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()

// 	pgx.MockerSet("Row.Scan", 1)
// 	err = AddOrderCommCall(Pool(), comms...)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Pool.Query", 1)
// 	_, _, err = ListUserOrderComm(100, order.TID)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Rows.Scan", 1)
// 	_, _, err = ListUserOrderComm(100, order.TID)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Pool.Query", 1)
// 	_, err = CountUserOrderComm(100, time.Now().Add(-10*time.Hour), time.Now())
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Rows.Scan", 1)
// 	_, err = CountUserOrderComm(100, time.Now().Add(-10*time.Hour), time.Now())
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Pool.Query", 1)
// 	_, err = CountMyUserOrderComm(100, 10)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Rows.Scan", 1)
// 	_, err = CountMyUserOrderComm(100, 10)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// }
