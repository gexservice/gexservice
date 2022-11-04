package gexdb

// import (
// 	"testing"
// 	"time"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/shopspring/decimal"
// )

// func TestProcSettleOrderFee(t *testing.T) {
// 	clear()
// 	user := testAddUser("TestProcSettleOrderFee-User")
// 	broker := testAddUser("TestProcSettleOrderFee-Broker")
// 	user.BrokerID = broker.TID
// 	UpdateUser(user)
// 	TouchBalance(BalanceAssetAll, user.TID, broker.TID)
// 	//
// 	order := &Order{
// 		Type:       OrderTypeTrade,
// 		UserID:     user.TID,
// 		FeeBalance: BalanceAssetYWE,
// 		FeeFilled:  decimal.NewFromFloat(10),
// 		Status:     OrderStatusDone,
// 	}
// 	err := CreateOrder(order)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ProcSettleOrderFee()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	assetBalanceFree(broker.TID, BalanceAssetYWE, decimal.NewFromFloat(8))

// 	comms, _, err := ListUserOrderComm(broker.TID, order.TID)
// 	if err != nil || len(comms) < 1 {
// 		t.Error("error")
// 		return
// 	}

// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	order = &Order{
// 		Type:       OrderTypeTrade,
// 		UserID:     user.TID,
// 		FeeBalance: BalanceAssetYWE,
// 		FeeFilled:  decimal.NewFromFloat(10),
// 		Status:     OrderStatusDone,
// 	}
// 	err = CreateOrder(order)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	for i := int64(1); i <= 6; i++ {
// 		Pool().Exec(`update gex_order set fee_settled_next=$1 where tid=$2`, time.Time{}, order.TID)
// 		pgx.MockerClear()
// 		pgx.MockerSet("Row.Scan", i)
// 		err = ProcSettleOrderFee()
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 	}

// 	Pool().Exec(`update gex_order set fee_settled_next=$1 where tid=$2`, time.Time{}, order.TID)
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	err = ProcSettleOrderFee()
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// }
