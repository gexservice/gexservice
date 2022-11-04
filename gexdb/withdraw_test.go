package gexdb

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/codingeasygo/util/xsql"
// 	"github.com/shopspring/decimal"
// )

// func TestWithdraw(t *testing.T) {
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		WithdrawVerifyOrder("")
// 	}()
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		WithdrawApplyOrder(nil, nil)
// 	}()
// 	clear()
// 	userMMK := testAddUser("TestMatcherMarket-MMK")
// 	userNone := testAddUser("TestMatcherMarket-NONE")
// 	_, err := TouchBalance(BalanceAssetAll, userMMK.TID, userNone.TID)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	IncreaseBalanceCall(Pool(), &Balance{
// 		UserID: userMMK.TID,
// 		Asset:  BalanceAssetMMK,
// 		Free:   decimal.NewFromFloat(10000),
// 		Status: BalanceStatusNormal,
// 	})
// 	{ //create>cancel
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		order, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err != nil || order.Status != OrderStatusCanceled {
// 			t.Error(err)
// 			return
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusCanceled {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply>done
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		err = ProcWithdrawApply()
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply>query>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyPending
// 			info = xmap.M{"pending": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 2; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			err = ProcWithdrawApply()
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply>fail>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		applyCount := 1
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			if applyCount > 0 {
// 				result = WithdrawVerifyFail
// 				info = xmap.M{"fail": 1}
// 				applyCount--
// 			} else {
// 				result = WithdrawVerifyPending
// 				info = xmap.M{"pending": 1}
// 			}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 3; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			err = ProcWithdrawApply()
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply error>verify not>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"message": "not found"}
// 			return
// 		}
// 		applyCount := 1
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			if applyCount > 0 {
// 				err = fmt.Errorf("error")
// 				applyCount--
// 			} else {
// 				result = WithdrawVerifyDone
// 				info = xmap.M{"done": 1}
// 			}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 3; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			err = ProcWithdrawApply()
// 			if i == 0 && err == nil {
// 				t.Error(err)
// 			}
// 			if i > 0 && err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply unknow>verify not>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"message": "not found"}
// 			return
// 		}
// 		applyCount := 1
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			if applyCount > 0 {
// 				result = -1
// 				info = xmap.M{"unknow": 1}
// 				applyCount--
// 			} else {
// 				result = WithdrawVerifyDone
// 				info = xmap.M{"done": 1}
// 			}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 3; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			err = ProcWithdrawApply()
// 			if i == 0 && err == nil {
// 				t.Error(err)
// 			}
// 			if i > 0 && err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>apply fail>cancel
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		err = ProcWithdrawApply()
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		order, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err != nil || order.Status != OrderStatusCanceled {
// 			t.Error(err)
// 			return
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusCanceled {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	{ //create>find error>verify not>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"message": "not found"}
// 			return
// 		}
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 4; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			pgx.MockerClear()
// 			if i < 2 {
// 				pgx.MockerSet("Row.Scan", int64(i+2))
// 			}
// 			err = ProcWithdrawApply()
// 			if i == 0 && err == nil {
// 				t.Error(err)
// 			}
// 			pgx.MockerClear()
// 			if i > 1 && err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>query error>verify not>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"message": "not found"}
// 			return
// 		}
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 3; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			pgx.MockerClear()
// 			if i == 0 {
// 				pgx.MockerSet("Row.Scan", 1)
// 			}
// 			err = ProcWithdrawApply()
// 			pgx.MockerClear()
// 			if i == 0 && err == nil {
// 				t.Error(err)
// 			}
// 			pgx.MockerClear()
// 			if i > 0 && err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //create>panic>verify not>apply>done
// 		WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyFail
// 			info = xmap.M{"message": "not found"}
// 			return
// 		}
// 		WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 			result = WithdrawVerifyDone
// 			info = xmap.M{"done": 1}
// 			return
// 		}
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil || order.Status != OrderStatusPending {
// 			t.Error(err)
// 			return
// 		}
// 		for i := 0; i < 3; i++ {
// 			err = Pool().ExecRow(`update gex_order set withdraw_next=$1 where tid=$2`, xsql.Time{}, order.TID)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			pgx.MockerClear()
// 			if i == 0 {
// 				pgx.MockerPanic("Row.Scan", 1)
// 			}
// 			err = ProcWithdrawApply()
// 			pgx.MockerClear()
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 		}
// 		order, err = FindOrderByOrderID(order.OrderID)
// 		if err != nil || order.Status != OrderStatusDone {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //db error
// 		//create error
// 		pgx.MockerSet("Pool.Begin", 1)
// 		_, err = CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		pgx.MockerSet("Tx.Exec", 1)
// 		_, err = CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		order, err := CreateWithdrawOrder(userMMK.TID, decimal.NewFromFloat(100))
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		//cancel error
// 		pgx.MockerSet("Pool.Begin", 1)
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		pgx.MockerSet("Row.Scan", 1)
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		pgx.MockerSet("Tx.Exec", 1)
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		_, err = CancelWithdrawOrder(10, order.OrderID) //not access
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID) //status error
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()

// 		//type error
// 		order = &Order{
// 			Type:      OrderTypeTopup,
// 			UserID:    userMMK.TID,
// 			InBalance: BalanceAssetMMK,
// 			InFilled:  decimal.NewFromFloat(100),
// 			Status:    OrderStatusPending,
// 		}
// 		err = CreateOrder(order)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		pgx.MockerClear()
// 		_, err = CancelWithdrawOrder(userMMK.TID, order.OrderID)
// 		if err == nil {
// 			t.Error(err)
// 			return
// 		}
// 	}
// }
